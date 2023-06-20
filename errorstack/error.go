package errorstack

import (
	"strings"
)

/*
Ensure *Error complies to Go error type.
*/
var _ error = (*Error)(nil)

/*
Error implements the Go native error type and is designed for handling errors
in the helix.go ecosystem. When exposing errors to clients (such as via HTTP
API), the root error should not give away too much information such as internal
messages.
*/
type Error struct {

	// Integration is the name of the integration returning the error, if applicable.
	// Omit integration when working with JSON: we don't want to give internal
	// information to clients consuming HTTP APIs.
	//
	// Examples:
	//
	//   "nats"
	//   "vault"
	Integration string `json:"-"`

	// Message is the top-level message of the error.
	Message string `json:"message"`

	// Validations represents a list of failed configuration validations. This is
	// used when a Integration's configuration encountered errors related to values
	// set by clients.
	Validations []Validation `json:"validations,omitempty"`

	// Children holds child errors encountered in cascade related to the current
	// error. Omit children errors when working with JSON: we don't want to give
	// internal information to clients consuming HTTP APIs.
	Children []error `json:"-"`
}

/*
Validation holds some details about a validation failure.
*/
type Validation struct {

	// Message is the cause of the validation failure.
	Message string `json:"message"`

	// Path represents the path to the key where the validation failure occurred.
	//
	// Example:
	//
	//   []string{"Options", "Producer", "MaxRetries"}
	Path []string `json:"path,omitempty"`
}

/*
New returns a new error given the message and options passed.
*/
func New(message string, opts ...With) *Error {
	err := &Error{
		Message:     message,
		Validations: []Validation{},
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

/*
NewFromError returns a new error given the existing error and options passed.
*/
func NewFromError(existing error, opts ...With) *Error {
	if existing == nil {
		return nil
	}

	err := &Error{
		Message:     existing.Error(),
		Validations: []Validation{},
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

/*
WithValidations adds validation failures to an error.
*/
func (err *Error) WithValidations(validations ...Validation) *Error {
	err.Validations = append(err.Validations, validations...)
	return err
}

/*
HasValidations indicates if an error encountered validation failures.
*/
func (err *Error) HasValidations() bool {
	return len(err.Validations) > 0
}

/*
WithChildren adds a list of child errors encountered related to the current
error.
*/
func (err *Error) WithChildren(children ...error) error {
	err.Children = append(err.Children, children...)
	return err
}

/*
HasChildren indicates if an error caused other (a.k.a. children) errors.
*/
func (err *Error) HasChildren() bool {
	return len(err.Children) > 0
}

/*
Error returns the stringified version of the error, including its validation
failures.
*/
func (err *Error) Error() string {
	var msg string

	if err.Integration != "" {
		msg += err.Integration + ": "
	}

	if err.Message != "" {
		msg += err.Message
	}

	if err.HasValidations() {
		msg += ". Reasons:\n"

		for _, validation := range err.Validations {
			msg += "    - " + validation.Message
			if validation.Path != nil && len(validation.Path) > 0 {
				msg += "\n      at " + strings.Join(validation.Path, " > ")
			}

			msg += ".\n"
		}
	} else {
		msg += "."
	}

	if err.HasChildren() {
		msg += " Caused by:"
		msg += "\n\n"

		for _, child := range err.Children {
			msg += "- " + child.Error()
			msg += "\n"
		}
	}

	return msg
}
