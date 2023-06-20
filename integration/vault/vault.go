package vault

import (
	"context"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"

	"github.com/hashicorp/vault/api"
)

/*
Vault exposes an opinionated way to interact with Vault, by bringing automatic
distributed tracing as well as error recording within traces.
*/
type Vault interface {
	KeyValue(ctx context.Context, path string) KeyValue
}

/*
connection represents the vault integration. It respects the integration.Integration
and Vault interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new Vault client.
	config *Config

	// client is the connection made with the Vault server.
	client *api.Client
}

/*
Connect tries to connect to the Vault server given the Config. Returns an error
if Config is not valid or if the connection failed.
*/
func Connect(cfg Config) (Vault, error) {

	// No need to continue if Config is not valid.
	err := cfg.sanitize()
	if err != nil {
		return nil, err
	}

	// Start to build an error stack, so we can add validations as we go.
	stack := errorstack.New("Failed to initialize integration", errorstack.WithIntegration(identifier))
	conn := &connection{
		config: &cfg,
	}

	// Set the default Vault config.
	var opts = &api.Config{
		Address:      cfg.Address,
		AgentAddress: cfg.AgentAddress,
	}

	// Set TLS options only if enabled in Config.
	if cfg.TLS.Enabled {
		tls := &api.TLSConfig{
			CACert:        cfg.TLS.RootCAFiles[0],
			ClientCert:    cfg.TLS.CertFile,
			ClientKey:     cfg.TLS.KeyFile,
			TLSServerName: cfg.TLS.ServerName,
		}

		err := opts.ConfigureTLS(tls)
		if err != nil {
			stack.WithValidations(errorstack.Validation{
				Message: normalizeErrorMessage(err),
			})
		}
	}

	// Try to connect to the Vault server.
	conn.client, err = api.NewClient(opts)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})
	}

	// Override the namespace if necessary.
	if cfg.Namespace != "" {
		conn.client.SetNamespace(cfg.Namespace)
	}

	// Override the token if necessary.
	if cfg.Token != "" {
		conn.client.SetToken(cfg.Token)
	}

	// Stop here if error validations were encountered.
	if stack.HasValidations() {
		return nil, stack
	}

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

/*
KeyValue will lookup and bind to an existing Key-Value v2 at the mount path.
*/
func (conn *connection) KeyValue(ctx context.Context, mountpath string) KeyValue {
	store := &keyvalue{
		config:    conn.config,
		mountpath: mountpath,
		client:    conn.client.KVv2(mountpath),
	}

	return store
}
