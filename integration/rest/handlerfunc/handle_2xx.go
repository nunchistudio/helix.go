package handlerfunc

import (
	"net/http"

	"go.nunchi.studio/helix/integration/rest"
)

/*
Accepted writes a 202 status code and body to the HTTP response writer.
*/
func Accepted(rw http.ResponseWriter, req *http.Request) {
	res := &rest.Response{
		Status: http.StatusText(http.StatusAccepted),
	}

	writeResponse(http.StatusAccepted, rw, res)
}
