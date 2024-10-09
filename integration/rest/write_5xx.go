package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
)

/*
WriteInternalServerError writes a 500 status code and body to the HTTP response
writer.

WithData has no effect here since no data must be returned in 5xx responses.

Default error message (in English) is:

	We have been notified of this unexpected internal error
*/
func WriteInternalServerError(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusInternalServerError]),
	}

	writeResponse(http.StatusInternalServerError, rw, res, opts...)
}

/*
WriteServiceUnavailable writes a 503 status code and body to the HTTP response
writer.

WithData has no effect here since no data must be returned in 5xx responses.

Default error message (in English) is:

	Please try again in a few moments
*/
func WriteServiceUnavailable(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusServiceUnavailable),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusServiceUnavailable]),
	}

	writeResponse(http.StatusServiceUnavailable, rw, res, opts...)
}
