package rest

import (
	"net/http"
)

/*
WriteOK writes a 200 status code and body to the HTTP response writer.

WithErrorMessage and WithValidations have no effect here since no error is
returned in 2xx responses.
*/
func WriteOK(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusOK),
	}

	writeResponse(http.StatusOK, rw, res, opts...)
}

/*
WriteCreated writes a 201 status code and body to the HTTP response writer.

WithErrorMessage and WithValidations have no effect here since no error is
returned in 2xx responses.
*/
func WriteCreated(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusCreated),
	}

	writeResponse(http.StatusCreated, rw, res, opts...)
}

/*
WriteAccepted writes a 202 status code and body to the HTTP response writer.

WithErrorMessage and WithValidations have no effect here since no error is
returned in 2xx responses.
*/
func WriteAccepted(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusAccepted),
	}

	writeResponse(http.StatusAccepted, rw, res, opts...)
}
