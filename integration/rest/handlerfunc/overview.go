/*
Package handlerfunc provides consistent HTTP responses for most common use cases.
The calling function should return as soon as the handler function is called to
avoid writing multiple times to the response writer.

Example:

	MyHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	  if req.Header.Get("X-API-KEY") == "" {
	    handlerfunc.Unauthorized(rw, req)
	    return
	  }
	}
*/
package handlerfunc
