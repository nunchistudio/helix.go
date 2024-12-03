package integration

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"go.nunchi.studio/helix/errorstack"
)

/*
ConfigTLS is the common configuration for TLS across all integrations.
*/
type ConfigTLS struct {

	// Enabled enables TLS for the integration. When disabled, other fields are
	// ignored and can be empty.
	Enabled bool `json:"enabled"`

	// ServerName is used to verify the hostname on the returned certificates. It
	// is also included in the client's handshake to support virtual hosting unless
	// it is an IP address.
	ServerName string `json:"server_name,omitempty"`

	// InsecureSkipVerify controls whether a client verifies the server's certificate
	// chain and host name. If InsecureSkipVerify is true, crypto/tls accepts any
	// certificate presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to machine-in-the-middle attacks unless
	// custom verification is used.
	InsecureSkipVerify bool `json:"insecure_skip_verify"`

	// CertFile is the relative or absolute path to the certificate file.
	//
	// Example:
	//
	//   "./server.crt"
	CertFile string `json:"-"`

	// KeyFile is the relative or absolute path to the private key file.
	//
	// Example:
	//
	//   "./server.key"
	KeyFile string `json:"-"`

	// RootCAFiles allows to provide the RootCAs pool from a list of filenames.
	// This is not required by all integrations.
	RootCAFiles []string `json:"-"`
}

/*
Sanitize sets default values - if applicable - and validates the configuration.
Returns validation errors if configuration is not valid. This doesn't return a
standard error since this function shall only be called by integrations. This
allows to easily add error validations to an existing errorstack:

	stack.WithValidations(cfg.TLS.Sanitize()...)
*/
func (cfg *ConfigTLS) Sanitize() []errorstack.Validation {
	var validations []errorstack.Validation
	if !cfg.Enabled {
		return validations
	}

	if cfg.CertFile == "" {
		validations = append(validations, errorstack.Validation{
			Message: "CertFile must be set and not be empty",
			Path:    []string{"Config", "TLS", "CertFile"},
		})
	}

	if cfg.KeyFile == "" {
		validations = append(validations, errorstack.Validation{
			Message: "KeyFile must be set and not be empty",
			Path:    []string{"Config", "TLS", "KeyFile"},
		})
	}

	return validations
}

/*
ToStandardTLS tries to return a Go standard *tls.Config. Returns validation errors
if configuration is not valid. This doesn't return a standard error since this
function shall only be called by integrations. This allows to easily add error
validations to an existing errorstack.
*/
func (cfg *ConfigTLS) ToStandardTLS() (*tls.Config, []errorstack.Validation) {
	var validations []errorstack.Validation
	if !cfg.Enabled {
		return nil, validations
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		validations = append(validations, errorstack.Validation{
			Message: err.Error(),
		})
	}

	tlsConfig := &tls.Config{
		ServerName:         cfg.ServerName,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Certificates:       []tls.Certificate{cert},
	}

	if len(cfg.RootCAFiles) == 0 {
		return tlsConfig, nil
	}

	caCertPool := x509.NewCertPool()
	for _, ca := range cfg.RootCAFiles {
		caCert, err := os.ReadFile(ca)
		if err != nil {
			validations = append(validations, errorstack.Validation{
				Message: err.Error(),
			})

			continue
		}

		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			validations = append(validations, errorstack.Validation{
				Message: "failed to append root certificate from pem",
			})
		}
	}

	if len(validations) > 0 {
		return nil, validations
	}

	tlsConfig.RootCAs = caCertPool
	return tlsConfig, nil
}
