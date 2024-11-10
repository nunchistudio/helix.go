package rest

import (
	"go.nunchi.studio/helix/errorstack"
)

/*
WithOnError allows to override error message, set validation errors, and fill the
"metadata" object when returning 4xx and 5xx HTTP responses.
*/
type WithOnError func(res *Response)

/*
WithErrorMessageOnError overrides the default translated error message. The
translation of the default error message relies on the HTTP cookie or
"Accept-Language" header. Using WithErrorMessage means it's up to the calling
client to handle translation of the error message on its side (if desired).
*/
func WithErrorMessageOnError(msg string) WithOnError {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Message = msg
		}
	}
}

/*
WithValidationsOnError adds validation errors to the main error.
*/
func WithValidationsOnError(validations []errorstack.Validation) WithOnError {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Validations = validations
		}
	}
}

/*
WithMetadataOnError set the "metadata" object of the response.
*/
func WithMetadataOnError[T any](metadata T) WithOnError {
	return func(res *Response) {
		if res != nil {
			res.Metadata = metadata
		}
	}
}

/*
WithOnEmptyError allows to override error message and set error validations when
returning 4xx and 5xx HTTP responses with no "metadata".
*/
type WithOnEmptyError func(res *Response)

/*
WithErrorMessageOnEmptyError overrides the default translated error message. The
translation of the default error message relies on the HTTP cookie or
"Accept-Language" header. Using WithErrorMessageOnEmptyError means it's up to
the calling client to handle translation of the error message on its side (if
desired).
*/
func WithErrorMessageOnEmptyError(msg string) WithOnEmptyError {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Message = msg
		}
	}
}

/*
WithValidationsOnEmptyError adds validation errors to the main error.
*/
func WithValidationsOnEmptyError(validations []errorstack.Validation) WithOnEmptyError {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Validations = validations
		}
	}
}

/*
WithOnSuccess allows to set "metadata" and "data" objects when returning 2xx HTTP
responses.
*/
type WithOnSuccess func(res *Response)

/*
WithMetadataOnSuccess set the "metadata" object of the response.
*/
func WithMetadataOnSuccess[T any](metadata T) WithOnSuccess {
	return func(res *Response) {
		if res != nil {
			res.Metadata = metadata
		}
	}
}

/*
WithDataOnSuccess set the "data" object of the response.
*/
func WithDataOnSuccess[T any](data T) WithOnSuccess {
	return func(res *Response) {
		if res != nil && res.Error == nil {
			res.Data = data
		}
	}
}
