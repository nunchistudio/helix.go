package rest

import (
	"net/http"
)

/*
WriteOK writes a 200 status code and body to the HTTP response writer.
*/
func WriteOK[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnSuccess) {
	res := &Response{
		Status: http.StatusText(http.StatusOK),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnSuccess[T](http.StatusOK, rw, res, req, opts...)
}

/*
WriteEmptyOK returns WriteOK with no "metadata" and no "data".
*/
func WriteEmptyOK(rw http.ResponseWriter, req *http.Request) {
	WriteOK[struct{}](rw, req)
}

/*
WriteCreated writes a 201 status code and body to the HTTP response writer.
*/
func WriteCreated[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnSuccess) {
	res := &Response{
		Status: http.StatusText(http.StatusCreated),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnSuccess[T](http.StatusCreated, rw, res, req, opts...)
}

/*
WriteEmptyCreated returns WriteCreated with no "metadata" and no "data".
*/
func WriteEmptyCreated(rw http.ResponseWriter, req *http.Request) {
	WriteCreated[struct{}](rw, req)
}

/*
WriteAccepted writes a 202 status code and body to the HTTP response writer.
*/
func WriteAccepted[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnSuccess) {
	res := &Response{
		Status: http.StatusText(http.StatusAccepted),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnSuccess[T](http.StatusAccepted, rw, res, req, opts...)
}

/*
WriteEmptyAccepted returns WriteAccepted with no "metadata" and no "data".
*/
func WriteEmptyAccepted(rw http.ResponseWriter, req *http.Request) {
	WriteAccepted[struct{}](rw, req)
}
