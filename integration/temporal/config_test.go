package temporal

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
				Namespace: "fake",
			},
			after: Config{
				Address:   "127.0.0.1:7233",
				Namespace: "fake",
			},
			err: nil,
		},
		{
			before: Config{},
			after: Config{
				Address:   "127.0.0.1:7233",
				Namespace: "default",
			},
			err: nil,
		},
		{
			before: Config{
				Worker: ConfigWorker{
					Enabled: true,
				},
			},
			after: Config{
				Address:   "127.0.0.1:7233",
				Namespace: "default",
				Worker: ConfigWorker{
					Enabled: true,
				},
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "TaskQueue must be set and not be empty",
						Path:    []string{"Config", "Worker", "TaskQueue"},
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
