package rest

import (
	"net/http"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
)

/*
Config is used to configure the HTTP REST integration.
*/
type Config struct {

	// Address is the HTTP address to listen on.
	//
	// Default:
	//
	//   ":8080"
	Address string `json:"address"`

	// Middleware allows to wrap the built-in HTTP handler with a custom one, for
	// adding a chain of middlewares.
	Middleware func(next http.Handler) http.Handler `json:"-"`

	// Healthcheck allows to define custom logic for the healthcheck endpoint at:
	//
	//   GET /health
	//
	// It should return 200 if service is healthy, or 5xx if an error occurred.
	// Returns 200 by default.
	Healthcheck func(req *http.Request) int `json:"-"`

	// OpenAPI configures OpenAPI behavior within the REST API.
	OpenAPI ConfigOpenAPI `json:"openapi"`

	// TLSConfig configures TLS for the HTTP server. Only CertFile and KeyFile
	// are took into consideration. Filenames containing a certificate and matching
	// private key for the server must be provided. If the certificate is signed
	// by a certificate authority, the CertFile should be the concatenation of the
	// server's certificate, any intermediates, and the CA's certificate.
	TLS integration.ConfigTLS `json:"tls"`
}

/*
ConfigOpenAPI configures OpenAPI behavior within the REST API. When enabled, HTTP
requests and responses are automatically validated againt the description passed.
If a request is not valid, a 4xx error is returned to the client. If a response
is not valid, an error is logged but the response is still returned to the client.
*/
type ConfigOpenAPI struct {

	// Enabled enables OpenAPI within the REST API.
	Enabled bool `json:"enabled"`

	// Description is a path to a local file or a URL containing the OpenAPI
	// description.
	//
	// Examples:
	//
	//   "./descriptions/openapi.yaml"
	//   "http://domain.tld/openapi.yaml"
	Description string `json:"description,omitempty"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if cfg.Address == "" {
		cfg.Address = ":8080"
	}

	if cfg.OpenAPI.Enabled {
		if cfg.OpenAPI.Description == "" {
			stack.WithValidations(errorstack.Validation{
				Message: "Description must be set and not be empty",
				Path:    []string{"Config", "OpenAPI", "Description"},
			})
		}
	}

	stack.WithValidations(cfg.TLS.Sanitize()...)
	if stack.HasValidations() {
		return stack
	}

	return nil
}
