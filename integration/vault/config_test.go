package vault

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
				Token: "fake",
			},
			after: Config{
				Address: "http://127.0.0.1:8200",
				Token:   "fake",
			},
			err: nil,
		},
		{
			before: Config{},
			after: Config{
				Address: "http://127.0.0.1:8200",
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Token must be set and not be empty",
						Path:    []string{"Config", "Token"},
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
