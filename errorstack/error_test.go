package errorstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testcases := []struct {
		input    *Error
		expected *Error
	}{
		{
			input: New("This is a simple text example"),
			expected: &Error{
				Message:     "This is a simple text example",
				Validations: []Validation{},
			},
		},
		{
			input: New("This is a text example with integration", WithIntegration("vault")),
			expected: &Error{
				Integration: "vault",
				Message:     "This is a text example with integration",
				Validations: []Validation{},
			},
		},
	}

	for _, tc := range testcases {
		actual := tc.input

		assert.Equal(t, tc.expected, actual)
	}
}

func TestError_Error(t *testing.T) {
	testcases := []struct {
		input    *Error
		expected string
	}{
		{
			input:    New("This is a simple text example"),
			expected: `This is a simple text example.`,
		},
		{
			input:    New("This is a text example with integration", WithIntegration("vault")),
			expected: `vault: This is a text example with integration.`,
		},
		{
			input: New("This is a text example with validations", WithIntegration("vault")).WithValidations(Validation{
				Message: "Failed to validate test case",
				Path:    []string{"custom", "path"},
			}),
			expected: `vault: This is a text example with validations. Reasons:
    - Failed to validate test case
      at custom > path.
`,
		},
	}

	for _, tc := range testcases {
		actual := tc.input.Error()

		assert.Equal(t, tc.expected, actual)
	}
}
