package bucket

import (
	"strings"

	"go.nunchi.studio/helix/errorstack"
)

/*
Config is used to configure the Bucket integration.
*/
type Config struct {

	// Driver sets the driver to use.
	//
	// Required.
	//
	// Example:
	//
	//   bucket.DriverAWS
	Driver Driver `json:"driver"`

	// Bucket is the name of the bucket.
	//
	// Required.
	//
	// Example:
	//
	//   "my-bucket"
	Bucket string `json:"bucket"`

	// Subfolder sets an optional subfolder where all keys are stored in the bucket.
	//
	// Default:
	//
	//   "/"
	//
	// Example:
	//
	//   "my/subfolder/"
	//
	// Operations on "<key>" will be translated to "my/subfolder/<key>".
	Subfolder string `json:"subfolder,omitempty"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if cfg.Driver == nil {
		stack.WithValidations(errorstack.Validation{
			Message: "Driver must be set and not be nil",
			Path:    []string{"Config", "Driver"},
		})
	} else {
		stack.WithValidations(cfg.Driver.validate(cfg)...)
	}

	if cfg.Bucket == "" {
		stack.WithValidations(errorstack.Validation{
			Message: "Bucket must be set and not be empty",
			Path:    []string{"Config", "Bucket"},
		})
	}

	if cfg.Subfolder != "" {
		if !strings.HasSuffix(cfg.Subfolder, "/") {
			stack.WithValidations(errorstack.Validation{
				Message: "Subfolder must end with a trailing slash",
				Path:    []string{"Config", "Subfolder"},
			})
		}
	}

	if stack.HasValidations() {
		return stack
	}

	return nil
}
