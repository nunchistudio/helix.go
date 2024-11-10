package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
)

/*
WriteBadRequest writes a 400 status code and body to the HTTP response writer.

Default error message (in English) is:

	Failed to validate request
*/
func WriteBadRequest[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusBadRequest]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusBadRequest, rw, res, req, opts...)
}

/*
WriteEmptyBadRequest returns WriteBadRequest with no "metadata".
*/
func WriteEmptyBadRequest(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteBadRequest[struct{}](rw, req)
}

/*
WriteUnauthorized writes a 401 status code and body to the HTTP response writer.

Default error message (in English) is:

	You are not authorized to perform this action
*/
func WriteUnauthorized[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusUnauthorized),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusUnauthorized]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusUnauthorized, rw, res, req, opts...)
}

/*
WriteEmptyUnauthorized returns WriteUnauthorized with no "metadata".
*/
func WriteEmptyUnauthorized(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteUnauthorized[struct{}](rw, req)
}

/*
WritePaymentRequired writes a 402 status code and body to the HTTP response writer.

Default error message (in English) is:

	Request failed because payment is required
*/
func WritePaymentRequired[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusPaymentRequired),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusPaymentRequired]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusPaymentRequired, rw, res, req, opts...)
}

/*
WriteEmptyPaymentRequired returns WritePaymentRequired with no "metadata".
*/
func WriteEmptyPaymentRequired(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WritePaymentRequired[struct{}](rw, req)
}

/*
WriteForbidden writes a 403 status code and body to the HTTP response writer.

Default error message (in English) is:

	You don't have required permissions to perform this action
*/
func WriteForbidden[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusForbidden),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusForbidden]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusForbidden, rw, res, req, opts...)
}

/*
WriteEmptyForbidden returns WriteForbidden with no "metadata".
*/
func WriteEmptyForbidden(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteForbidden[struct{}](rw, req)
}

/*
WriteNotFound writes a 404 status code and body to the HTTP response writer.

Default error message (in English) is:

	Resource does not exist
*/
func WriteNotFound[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusNotFound),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusNotFound]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusNotFound, rw, res, req, opts...)
}

/*
WriteEmptyNotFound returns WriteNotFound with no "metadata".
*/
func WriteEmptyNotFound(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteNotFound[struct{}](rw, req)
}

/*
WriteMethodNotAllowed writes a 405 status code and body to the HTTP response writer.

Default error message (in English) is:

	Resource does not support this method
*/
func WriteMethodNotAllowed[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusMethodNotAllowed),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusMethodNotAllowed]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusMethodNotAllowed, rw, res, req, opts...)
}

/*
WriteEmptyMethodNotAllowed returns WriteMethodNotAllowed with no "metadata".
*/
func WriteEmptyMethodNotAllowed(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteMethodNotAllowed[struct{}](rw, req)
}

/*
WriteConflict writes a 409 status code and body to the HTTP response writer.

Default error message (in English) is:

	Failed to process target resource because of conflict
*/
func WriteConflict[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusConflict),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusConflict]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusConflict, rw, res, req, opts...)
}

/*
WriteEmptyConflict returns WriteConflict with no "metadata".
*/
func WriteEmptyConflict(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteConflict[struct{}](rw, req)
}

/*
WriteRequestEntityTooLarge writes a 413 status code and body to the HTTP response
writer.

Default error message (in English) is:

	Can not process payload too large
*/
func WriteRequestEntityTooLarge[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusRequestEntityTooLarge),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusRequestEntityTooLarge]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusRequestEntityTooLarge, rw, res, req, opts...)
}

/*
WriteEmptyRequestEntityTooLarge returns WriteRequestEntityTooLarge with no
"metadata".
*/
func WriteEmptyRequestEntityTooLarge(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteRequestEntityTooLarge[struct{}](rw, req)
}

/*
WriteTooManyRequests writes a 429 status code and body to the HTTP response writer.

Default error message (in English) is:

	Request-rate limit has been reached
*/
func WriteTooManyRequests[T any](rw http.ResponseWriter, req *http.Request, opts ...WithOnError) {
	res := &Response{
		Status: http.StatusText(http.StatusTooManyRequests),
		Error:  errorstack.New(supportedLocales[getPreferredLanguage(req)][http.StatusTooManyRequests]),
	}

	for _, opt := range opts {
		opt(res)
	}

	writeResponseOnError[T](http.StatusTooManyRequests, rw, res, req, opts...)
}

/*
WriteEmptyTooManyRequests returns WriteTooManyRequests with no "metadata".
*/
func WriteEmptyTooManyRequests(rw http.ResponseWriter, req *http.Request, opts ...WithOnEmptyError) {
	WriteTooManyRequests[struct{}](rw, req)
}
