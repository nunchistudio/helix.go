package openfeature

import (
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Sanitize(t *testing.T) {
	testcases := []struct {
		before Config
		after  Config
		err    error
	}{
		{
			before: Config{
				Paths: []string{"./example.yaml"},
			},
			after: Config{
				Paths: []string{"./example.yaml"},
			},
			err: nil,
		},
		{
			before: Config{
				Paths: nil,
			},
			after: Config{
				Paths: nil,
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Paths must be set and not be empty",
						Path:    []string{"Config", "Paths"},
					},
				},
			},
		},
		{
			before: Config{
				Paths: []string{},
			},
			after: Config{
				Paths: []string{},
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Paths must be set and not be empty",
						Path:    []string{"Config", "Paths"},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		err := tc.before.sanitize()

		assert.Equal(t, tc.before, tc.after)
		assert.Equal(t, tc.err, err)
	}
}
