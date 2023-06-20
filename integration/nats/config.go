package nats

import (
	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
)

/*
Config is used to configure the NATS integration.
*/
type Config struct {

	// Addresses are NATS addresses to connect to. A URL can contain username and
	// password, such as:
	//
	//   "nats://derek:pass@localhost:4222"
	//
	// Default:
	//
	//   []string{"nats://localhost:4222"}
	Addresses []string `json:"addresses"`

	// TLSConfig configures TLS to communicate with the NATS server.
	TLS integration.ConfigTLS `json:"tls"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if len(cfg.Addresses) == 0 {
		cfg.Addresses = []string{"nats://localhost:4222"}
	}

	stack.WithValidations(cfg.TLS.Sanitize()...)
	if stack.HasValidations() {
		return stack
	}

	return nil
}
