package openfeature

import (
	"go.nunchi.studio/helix/errorstack"
)

/*
Config is used to configure the OpenFeature integration.
*/
type Config struct {

	// Paths are the file paths containing the GO Feature Flag strategies.
	//
	// Required.
	Paths []string `json:"paths"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if cfg.Paths == nil || len(cfg.Paths) == 0 {
		stack.WithValidations(errorstack.Validation{
			Message: "Paths must be set and not be empty",
			Path:    []string{"Config", "Paths"},
		})
	}

	if stack.HasValidations() {
		return stack
	}

	return nil
}
