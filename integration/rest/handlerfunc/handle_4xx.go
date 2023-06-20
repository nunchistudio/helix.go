package handlerfunc

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration/rest"
)

/*
BadRequest writes a 400 status code and body to the HTTP response writer.
Additional error validations can be passed to give clients some details about
the failure.
*/
func BadRequest(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  errorstack.New("Failed to validate request"),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponse(http.StatusBadRequest, rw, res)
}

/*
Unauthorized writes a 401 status code and body to the HTTP response writer.
Additional error validations can be passed to give clients some details about
the failure.
*/
func Unauthorized(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusUnauthorized),
		Error:  errorstack.New("You are not authorized to perform this action"),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponse(http.StatusUnauthorized, rw, res)
}

/*
Forbidden writes a 403 status code and body to the HTTP response writer.
Additional error validations can be passed to give clients some details about
the failure.
*/
func Forbidden(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusForbidden),
		Error:  errorstack.New("You don't have required permissions to perform this action"),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponse(http.StatusForbidden, rw, res)
}

/*
NotFound writes a 404 status code and body to the HTTP response writer.
*/
func NotFound(rw http.ResponseWriter, req *http.Request) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusNotFound),
		Error:  errorstack.New("Resource does not exist"),
	}

	writeResponse(http.StatusNotFound, rw, res)
}

/*
MethodNotAllowed writes a 404 status code and body to the HTTP response writer.
*/
func MethodNotAllowed(rw http.ResponseWriter, req *http.Request) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusMethodNotAllowed),
		Error:  errorstack.New("Resource does not support this method"),
	}

	writeResponse(http.StatusMethodNotAllowed, rw, res)
}

/*
RequestEntityTooLarge writes a 413 status code and body to the HTTP response
writer. Additional error validations can be passed to give clients some details
about the failure.
*/
func RequestEntityTooLarge(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusRequestEntityTooLarge),
		Error:  errorstack.New("Can not process payload too large"),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponse(http.StatusRequestEntityTooLarge, rw, res)
}

/*
TooManyRequests writes a 429 status code and body to the HTTP response writer.
Additional error validations can be passed to give clients some details about
the failure.
*/
func TooManyRequests(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusTooManyRequests),
		Error:  errorstack.New("Request-rate limit has been reached"),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponse(http.StatusTooManyRequests, rw, res)
}
