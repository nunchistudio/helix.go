package rest

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
			before: Config{},
			after: Config{
				Address: ":8080",
			},
			err: nil,
		},
		{
			before: Config{
				OpenAPI: ConfigOpenAPI{
					Enabled: true,
				},
			},
			after: Config{
				Address: ":8080",
				OpenAPI: ConfigOpenAPI{
					Enabled: true,
				},
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Description must be set and not be empty",
						Path:    []string{"Config", "OpenAPI", "Description"},
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
