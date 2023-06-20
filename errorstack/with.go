package errorstack

/*
With allows to set optional values when creating a new error with New.
*/
type With func(*Error)

/*
WithIntegration sets the integration at the origin of the error.
*/
func WithIntegration(inte string) With {
	return func(err *Error) {
		err.Integration = inte
	}
}
