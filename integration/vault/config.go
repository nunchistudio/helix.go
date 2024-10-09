package vault

import (
	"os"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
)

/*
Config is used to configure the Vault integration.
*/
type Config struct {

	// Address is the Vault server address to connect to. This should be a complete
	// URL.
	//
	// Default:
	//
	//   "http://127.0.0.1:8200"
	Address string `json:"address"`

	// AgentAddress is the local Vault agent address to connect to. This should be
	// a complete URL.
	//
	// Example:
	//
	//   "http://127.0.0.1:8200"
	AgentAddress string `json:"agent_address"`

	// Namespace sets the namespace to connect to, if not already set via environment
	// variable.
	Namespace string `json:"namespace"`

	// Token sets the token to use, if not already set via environment variable.
	Token string `json:"-"`

	// TLSConfig configures TLS to communicate with the Vault server.
	//
	// Important: Only the first Root CA file will be used and applied.
	TLS integration.ConfigTLS `json:"tls"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if cfg.Address == "" {
		cfg.Address = "http://127.0.0.1:8200"
	}

	if cfg.Namespace == "" && os.Getenv("VAULT_NAMESPACE") != "" {
		cfg.Namespace = os.Getenv("VAULT_NAMESPACE")
	}

	if cfg.Token == "" && os.Getenv("VAULT_TOKEN") == "" {
		stack.WithValidations(errorstack.Validation{
			Message: "Token must be set and not be empty",
			Path:    []string{"Config", "Token"},
		})
	}

	stack.WithValidations(cfg.TLS.Sanitize()...)
	if stack.HasValidations() {
		return stack
	}

	return nil
}
