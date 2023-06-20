package handlerfunc

import (
	"encoding/json"
	"net/http"

	"go.nunchi.studio/helix/integration/rest"
)

/*
writeResponse is a small utility used by all HTTP handler functions of this
package to easily write response's status code and body.
*/
func writeResponse(status int, rw http.ResponseWriter, res *rest.Response) {
	b, _ := json.Marshal(res)
	rw.WriteHeader(status)
	rw.Write(b)
}
