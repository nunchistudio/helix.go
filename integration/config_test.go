package integration

import (
	"testing"

	"go.nunchi.studio/helix/errorstack"

	"github.com/stretchr/testify/assert"
)

func TestConfigTLS_Sanitize(t *testing.T) {
	testcases := []struct {
		before      ConfigTLS
		after       ConfigTLS
		validations []errorstack.Validation
	}{
		{
			before: ConfigTLS{},
			after: ConfigTLS{
				Enabled: false,
			},
			validations: nil,
		},
		{
			before: ConfigTLS{
				Enabled: true,
			},
			after: ConfigTLS{
				Enabled: true,
			},
			validations: []errorstack.Validation{
				{
					Message: "CertFile must be set and not be empty",
					Path:    []string{"Config", "TLS", "CertFile"},
				},
				{
					Message: "KeyFile must be set and not be empty",
					Path:    []string{"Config", "TLS", "KeyFile"},
				},
			},
		},
	}

	for _, tc := range testcases {
		validations := tc.before.Sanitize()

		assert.Equal(t, tc.before, tc.after)
		assert.Equal(t, tc.validations, validations)
	}
}
