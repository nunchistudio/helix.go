package bucket

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
				Driver: DriverLocal,
				Bucket: "anything",
			},
			after: Config{
				Driver: DriverLocal,
				Bucket: "anything",
			},
			err: nil,
		},
		{
			before: Config{
				Driver: DriverLocal,
			},
			after: Config{
				Driver: DriverLocal,
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Bucket must be set and not be empty",
						Path:    []string{"Config", "Bucket"},
					},
				},
			},
		},
		{
			before: Config{},
			after:  Config{},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Driver must be set and not be nil",
						Path:    []string{"Config", "Driver"},
					},
					{
						Message: "Bucket must be set and not be empty",
						Path:    []string{"Config", "Bucket"},
					},
				},
			},
		},
		{
			before: Config{
				Driver:    DriverLocal,
				Bucket:    "anything",
				Subfolder: "not/a/valid/path",
			},
			after: Config{
				Driver:    DriverLocal,
				Bucket:    "anything",
				Subfolder: "not/a/valid/path",
			},
			err: &errorstack.Error{
				Integration: identifier,
				Message:     "Failed to validate configuration",
				Validations: []errorstack.Validation{
					{
						Message: "Subfolder must end with a trailing slash",
						Path:    []string{"Config", "Subfolder"},
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
