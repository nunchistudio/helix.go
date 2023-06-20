package clickhouse

import (
	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
)

/*
Config is used to configure the ClickHouse integration.
*/
type Config struct {

	// Addresses are ClickHouse addresses to connect to.
	//
	// Default:
	//
	//   []string{"127.0.0.1:8123"}
	Addresses []string `json:"addresses"`

	// Database is the database to connect to.
	//
	// Default:
	//
	//   "default"
	Database string `json:"-"`

	// User is the user to use to connect to the database.
	User string `json:"-"`

	// Password is the user's password to connect to the database.
	Password string `json:"-"`

	// TLSConfig configures TLS to communicate with the ClickHouse server.
	TLS integration.ConfigTLS `json:"tls"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if len(cfg.Addresses) == 0 {
		cfg.Addresses = []string{"127.0.0.1:8123"}
	}

	if cfg.Database == "" {
		cfg.Database = "default"
	}

	stack.WithValidations(cfg.TLS.Sanitize()...)
	if stack.HasValidations() {
		return stack
	}

	return nil
}
