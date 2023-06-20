package handlerfunc

import (
	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration/rest"
)

/*
With allows to set optional values when calling some HTTP handler functions.
*/
type With func(res *rest.Response)

/*
WithValidations adds in-house validation errors to the main error.
*/
func WithValidations(validations []errorstack.Validation) With {
	return func(res *rest.Response) {
		if res != nil && res.Error != nil {
			res.Error.Validations = append(res.Error.Validations, validations...)
		}
	}
}
