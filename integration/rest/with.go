package rest

import (
	"go.nunchi.studio/helix/errorstack"
)

/*
With allows to set optional values when calling some HTTP handler functions.
*/
type With func(res *Response)

/*
WithErrorMessage overrides the default translated error message. The translation
of the default error message relies on the HTTP cookie or "Accept-Language" header.
Using WithErrorMessage means it's up to the calling client to handle translation
of the error message on its side (if desired).

WithErrorMessage has no effect on 2xx HTTP responses.
*/
func WithErrorMessage(msg string) With {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Message = msg
		}
	}
}

/*
WithValidations adds validation errors to the main error.

WithValidations has no effect on 2xx HTTP responses.
*/
func WithValidations(validations []errorstack.Validation) With {
	return func(res *Response) {
		if res != nil && res.Error != nil {
			res.Error.Validations = append(res.Error.Validations, validations...)
		}
	}
}

/*
WithMetadata set the metadata object for the response.
*/
func WithMetadata[T any](metadata T) With {
	return func(res *Response) {
		if res != nil {
			res.Metadata = metadata
		}
	}
}

/*
WithData set the data object for the response.

WithData has no effect on non-2xx HTTP responses.
*/
func WithData[T any](data T) With {
	return func(res *Response) {
		if res != nil && res.Error == nil {
			res.Data = data
		}
	}
}
