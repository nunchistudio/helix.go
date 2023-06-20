package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/uptrace/bunrouter"
)

/*
middlewareValidation is the HTTP middleware to validate a request/response against
the OpenAPI description passed in the integration's config.
*/
func (r *rest) middlewareValidation(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {

		// Create a new trace for the OpenAPI middleware. Since there's already a
		// trace in the request's context, spans will be part of the parent trace.
		ctx, spanReq := trace.Start(req.Context(), trace.SpanKindServer, "OpenAPI: Request validation")

		// Wrap the standard http.ResponseWriter so we can store additional values
		// during the request/response lifecycle, such as the status code and the
		// the response body.
		rw := &responseWriter{
			status:         200,
			ResponseWriter: w,
			buf:            &bytes.Buffer{},
		}

		// Try to find the route in the OpenAPI description. If the path is not found
		// or if the method is not allowed, it's already catched by the router itself
		// so there's no need to handle this here.
		r, params, err := r.oapirouter.FindRoute(req.Request)
		if err != nil {
			spanReq.RecordError("failed to find route", err)
			spanReq.End()
			return next(rw, req)
		}

		// Build the request input for OpenAPI validation. Only validate the
		// authentication if the security scheme is present and is in the headers
		// of the request.
		in := &openapi3filter.RequestValidationInput{
			Request:     req.Request,
			PathParams:  params,
			QueryParams: req.URL.Query(),
			Route:       r,
			Options: &openapi3filter.Options{
				MultiError: true,
				AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
					if ai != nil && ai.SecurityScheme != nil {
						if ai.SecurityScheme.In == "header" {
							if strings.TrimSpace(req.Header.Get(ai.SecurityScheme.Name)) == "" {
								return fmt.Errorf(ai.SecurityScheme.Name)
							}
						}
					}

					return nil
				},
			},
		}

		// Override the request validation error's string to only return the reason
		// of the error, with no additional details.
		in.Options.WithCustomSchemaErrorFunc(func(err *openapi3.SchemaError) string {
			return err.Reason
		})

		// Whatever happens next, make sure to validate the response returned, just
		// like we did for the request. If the response is not valid, an error is
		// recorded but the response is still returned to the client.
		defer func() {
			ctx, spanRes := trace.Start(req.Context(), trace.SpanKindServer, "OpenAPI: Response validation")

			out := &openapi3filter.ResponseValidationInput{
				RequestValidationInput: in,
				Status:                 rw.status,
				Header:                 rw.Header(),
				Body:                   io.NopCloser(rw.buf),
				Options: &openapi3filter.Options{
					MultiError:            true,
					IncludeResponseStatus: true,
				},
			}

			out.Options.WithCustomSchemaErrorFunc(func(err *openapi3.SchemaError) string {
				return err.Reason
			})

			err = openapi3filter.ValidateResponse(ctx, out)
			if err != nil {
				spanRes.RecordError("failed to validate response", err)
			}

			spanRes.End()
		}()

		// We now can validate the request. If the request does not respect the
		// OpenAPI description, return a 400 error and stop the request/response
		// lifecycle. We don't want HTTP handlers being called if the request is
		// not valid.
		err = openapi3filter.ValidateRequest(ctx, in)
		if err != nil {
			res := &Response{
				Status: "Bad Request",
				Error:  errorstack.New("Failed to validate request"),
			}

			// Convert the default error to match the errorsstack.Error format. Add
			// each "issue" to the slice of validations with their message and path.
			switch err := err.(type) {
			case openapi3.MultiError:
				issues := convertError("request.body", err)
				names := make([]string, 0, len(issues))
				for k := range issues {
					names = append(names, k)
				}

				sort.Strings(names)
				for _, k := range names {
					msgs := issues[k]
					for _, msg := range msgs {
						res.Error.Validations = append(res.Error.Validations, errorstack.Validation{
							Message: msg,
							Path:    strings.Split(k, "."),
						})
					}
				}
			}

			spanReq.RecordError("failed to validate request", err)
			spanReq.End()

			// Write the error validations to the response writer, informing the client
			// the errors encountered.
			b, _ := json.Marshal(res)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(b)

			return nil
		}

		// If we made it here it means the request is valid. We can close the span
		// and move to the next HTTP handler function.
		spanReq.End()
		next(rw, req)

		return nil
	}
}
