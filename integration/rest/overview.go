/*
Package rest exposes opinionated HTTP REST resources respecting the standards of
the net/http package. In addition to a strong association with OpenTelemetry, it
comes with OpenAPI support as well for request/response validation.

When using HTTP response writers, the calling function should return as soon as
the handler function is called to avoid writing multiple times to the response
writer.

Example:

	MyHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	  if req.Header.Get("X-API-KEY") == "" {
	    rest.Unauthorized(rw, req)
	    return
	  }
	}
*/
package rest

/*
identifier represents the integration's unique identifier.
*/
const identifier = "rest"

/*
humanized represents the integration's humanized name.
*/
// const humanized = "REST"
