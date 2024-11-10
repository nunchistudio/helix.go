package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
)

/*
WriteInternalServerError writes a 500 status code and body to the HTTP response
writer.

Default error message (in English) is:

	We have been notified of this unexpected internal error
*/
func WriteInternalServerError[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusInternalServerError]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusInternalServerError, rw, res, req, opts...)
}

/*
WriteEmptyInternalServerError returns WriteInternalServerError with no "metadata".
*/
func WriteEmptyInternalServerError(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteInternalServerError[struct{}](rw, req)
}

/*
WriteServiceUnavailable writes a 503 status code and body to the HTTP response
writer.

Default error message (in English) is:

	Please try again in a few moments
*/
func WriteServiceUnavailable[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusServiceUnavailable),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusServiceUnavailable]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusServiceUnavailable, rw, res, req, opts...)
}

/*
WriteEmptyServiceUnavailable returns WriteServiceUnavailable with no "metadata".
*/
func WriteEmptyServiceUnavailable(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteServiceUnavailable[struct{}](rw, req)
}
