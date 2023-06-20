package rest

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.nunchi.studio/helix/errorstack"

	"github.com/uptrace/bunrouter"
)

/*
responseWriter wraps the standard http.ResponseWriter so we can store additional
values during the request/response lifecycle, such as the status code and the
the response body.
*/
type responseWriter struct {
	http.ResponseWriter

	// status code is the HTTP status code sets in the response header. This allows
	// to ensure if the status code respects the one defined in the OpenAPI
	// description.
	status int

	// buf is the HTTP response body sets by a handler function. This allows to
	// ensure if the body respects the one defined in the OpenAPI description.
	buf *bytes.Buffer
}

/*
Write writes the data to the connection as part of an HTTP reply.
*/
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.ResponseWriter.Write(b)
	return rw.buf.Write(b)
}

/*
WriteHeader sends an HTTP response header with the provided status code.
*/
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

/*
Response is the JSON object every HTTP responses shall return.
*/
type Response struct {
	Status string            `json:"status"`
	Error  *errorstack.Error `json:"error,omitempty"`
	Data   any               `json:"data,omitempty"`
}

/*
handlerHealthcheck is the default handler function for the healthcheck endpoint.
Call the custom function defined in the Config if applicable.
*/
func (r *rest) handlerHealthcheck(rw http.ResponseWriter, req bunrouter.Request) error {
	var status int = http.StatusOK
	if r.config.Healthcheck != nil {
		status = r.config.Healthcheck(req.Request)
	}

	res := &Response{
		Status: http.StatusText(status),
	}

	b, _ := json.Marshal(res)
	rw.WriteHeader(status)
	rw.Write(b)

	return nil
}

/*
handlerNotFound is the default handler function if the path is not found (error
404).
*/
func (r *rest) handlerNotFound(rw http.ResponseWriter, req bunrouter.Request) error {
	res := &Response{
		Status: http.StatusText(http.StatusNotFound),
		Error: &errorstack.Error{
			Message: "Resource does not exist",
		},
	}

	b, _ := json.Marshal(res)
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(b)

	return nil
}

/*
handlerMethodNotAllowed is the default handler function if the method is not
allowed (error 405).
*/
func (r *rest) handlerMethodNotAllowed(rw http.ResponseWriter, req bunrouter.Request) error {
	res := &Response{
		Status: http.StatusText(http.StatusMethodNotAllowed),
		Error: &errorstack.Error{
			Message: "Resource does not support this method",
		},
	}

	b, _ := json.Marshal(res)
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(b)

	return nil
}
