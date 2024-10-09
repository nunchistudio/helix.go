package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
)

/*
WriteBadRequest writes a 400 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Failed to validate request
*/
func WriteBadRequest(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusBadRequest]),
	}

	writeResponse(http.StatusBadRequest, rw, res, opts...)
}

/*
WriteUnauthorized writes a 401 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	You are not authorized to perform this action
*/
func WriteUnauthorized(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusUnauthorized),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusUnauthorized]),
	}

	writeResponse(http.StatusUnauthorized, rw, res, opts...)
}

/*
WritePaymentRequired writes a 402 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Request failed because payment is required
*/
func WritePaymentRequired(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusPaymentRequired),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusPaymentRequired]),
	}

	writeResponse(http.StatusPaymentRequired, rw, res, opts...)
}

/*
WriteForbidden writes a 403 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	You don't have required permissions to perform this action
*/
func WriteForbidden(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusForbidden),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusForbidden]),
	}

	writeResponse(http.StatusForbidden, rw, res, opts...)
}

/*
WriteNotFound writes a 404 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Resource does not exist
*/
func WriteNotFound(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusNotFound),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusNotFound]),
	}

	writeResponse(http.StatusNotFound, rw, res, opts...)
}

/*
WriteMethodNotAllowed writes a 405 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Resource does not support this method
*/
func WriteMethodNotAllowed(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusMethodNotAllowed),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusMethodNotAllowed]),
	}

	writeResponse(http.StatusMethodNotAllowed, rw, res, opts...)
}

/*
WriteConflict writes a 409 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Failed to process target resource because of conflict
*/
func WriteConflict(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusConflict),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusConflict]),
	}

	writeResponse(http.StatusConflict, rw, res, opts...)
}

/*
WriteRequestEntityTooLarge writes a 413 status code and body to the HTTP response
writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Can not process payload too large
*/
func WriteRequestEntityTooLarge(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusRequestEntityTooLarge),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusRequestEntityTooLarge]),
	}

	writeResponse(http.StatusRequestEntityTooLarge, rw, res, opts...)
}

/*
WriteTooManyRequests writes a 429 status code and body to the HTTP response writer.

WithData has no effect here since no data must be returned in 4xx responses.

Default error message (in English) is:

	Request-rate limit has been reached
*/
func WriteTooManyRequests(rw http.ResponseWriter, req *http.Request, opts ...With) {
	res := &Response{
		Status: http.StatusText(http.StatusTooManyRequests),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusTooManyRequests]),
	}

	writeResponse(http.StatusTooManyRequests, rw, res, opts...)
}
