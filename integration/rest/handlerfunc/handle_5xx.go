package handlerfunc

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration/rest"
)

/*
InternalServerError writes a 500 status code and body to the HTTP response writer.
*/
func InternalServerError(rw http.ResponseWriter, req *http.Request) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Error:  errorstack.New("We have been notified of this unexpected internal error"),
	}

	writeResponse(http.StatusInternalServerError, rw, res)
}

/*
ServiceUnavailable writes a 503 status code and body to the HTTP response writer.
*/
func ServiceUnavailable(rw http.ResponseWriter, req *http.Request) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusServiceUnavailable),
		Error:  errorstack.New("Please try again in a few moments"),
	}

	writeResponse(http.StatusServiceUnavailable, rw, res)
}
