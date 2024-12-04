package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/telemetry/log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

/*
writeResponseOnError is a small utility used by 4xx and 5xx HTTP response writer
functions of this package to easily write response's status code and body.
*/
func writeResponseOnError[T any](status int, rw http.ResponseWriter, res *Response, req *http.Request, opts ...WithOnError) {
	for _, opt := range opts {
		opt(res)
	}

	writeResponse[T](status, rw, res, req)
}

/*
writeResponseOnSuccess is a small utility used by 2xx HTTP response writer
functions of this package to easily write response's status code and body.
*/
func writeResponseOnSuccess[T any](status int, rw http.ResponseWriter, res *Response, req *http.Request, opts ...WithOnSuccess) {
	for _, opt := range opts {
		opt(res)
	}

	writeResponse[T](status, rw, res, req)
}

/*
writeResponse is a small utility to easily write response's status code and body.
*/
func writeResponse[T any](status int, rw http.ResponseWriter, res *Response, req *http.Request) {
	b, err := json.Marshal(res)
	if err != nil {
		WriteEmptyInternalServerError(rw, req)
		return
	}

	var typed T
	if err := json.Unmarshal(b, &typed); err != nil {
		log.Error(req.Context(), "http response does not comply to struct `rest.Response`")
	}

	rw.WriteHeader(status)
	rw.Write(b)
}

/*
buildRouter tries to build the HTTP router. It comes with opinionated handlers
for 404 and 405 HTTP errors, as well as for the health endpoint.
*/
func (r *rest) buildRouter() (*bunrouter.CompatRouter, []errorstack.Validation) {
	opts := []bunrouter.Option{
		bunrouter.Use(reqlog.NewMiddleware(reqlog.WithEnabled(false))),
		bunrouter.Use(bunrouterotel.NewMiddleware(bunrouterotel.WithClientIP())),
		bunrouter.WithMiddleware(func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
			return func(rw http.ResponseWriter, req bunrouter.Request) error {
				rw.Header().Set("Content-Type", "application/json")
				return next(rw, req)
			}
		}),
		bunrouter.WithNotFoundHandler(r.handlerNotFound),
		bunrouter.WithMethodNotAllowedHandler(r.handlerMethodNotAllowed),
	}

	if r.config.OpenAPI.Enabled {
		opts = append(opts, bunrouter.WithMiddleware(r.middlewareValidation))
	}

	router := bunrouter.New(opts...).Compat()
	router.Router.GET("/health", r.handlerHealthcheck)

	return router, nil
}

/*
buildRouterOpenAPI tries to build the router for validating requests and responses
against the OpenAPI description. It returns validation errors in case the
description can not be loaded or if it's not valid.
*/
func (r *rest) buildRouterOpenAPI() (routers.Router, []errorstack.Validation) {
	loader := openapi3.NewLoader()

	// Load the description from file or from a URL, depending on the path defined
	// in the Config.
	var doc *openapi3.T
	var err error
	u, ok := isValidUrl(r.config.OpenAPI.Description)
	if ok {
		doc, err = loader.LoadFromURI(u)
	} else {
		doc, err = loader.LoadFromFile(r.config.OpenAPI.Description)
	}

	if err != nil {
		return nil, []errorstack.Validation{
			{
				Message: err.Error(),
			},
		}
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, []errorstack.Validation{
			{
				Message: err.Error(),
			},
		}
	}

	return router, nil
}

/*
isValidUrl tests a string to determine if it is a well-structured URL or not.
*/
func isValidUrl(link string) (*url.URL, bool) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return nil, false
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, false
	}

	return u, true
}

/*
convertError converts errors encountered by OpenAPI validation to a map of issues
so we can more easily format them in the response returned to the client.

Example:

	map[string][]string{
	  "request.headers.X-API-KEY": ["security requirements must be set"]
	  "request.body.event.name": ["property "name" is missing"]
	}
*/
func convertError(prefix string, me openapi3.MultiError) map[string][]string {
	issues := make(map[string][]string)

	for _, err := range me {
		switch err := err.(type) {

		// Check schema error.
		case *openapi3.SchemaError:
			field := prefix
			if path := err.JSONPointer(); len(path) > 0 {
				field = fmt.Sprintf("%s.%s", field, strings.Join(path, "."))
			}

			issues[field] = append(issues[field], err.Error())

		// Check global security requirements errors.
		case *openapi3filter.SecurityRequirementsError:
			field := fmt.Sprintf("%s.%s", "request.headers", strings.TrimPrefix(err.Error(), "security requirements failed: "))
			issues[field] = append(issues[field], "security requirements must be set")

		// Check request schema error.
		case *openapi3filter.RequestError:
			if err.Parameter != nil {
				field := fmt.Sprintf("%s.%s", err.Parameter.In, err.Parameter.Name)
				issues[field] = append(issues[field], err.Error())
				continue
			}

			if err, ok := err.Err.(openapi3.MultiError); ok {
				for k, v := range convertError(prefix, err) {
					issues[k] = append(issues[k], v...)
				}

				continue
			}

			if err.RequestBody != nil {
				issues[prefix] = append(issues[prefix], err.Error())
				continue
			}

		// Make sure to handle every usecases. Even though this should never happen.
		default:
			field := "unknown"
			issues[field] = append(issues[field], err.Error())
		}
	}

	return issues
}
