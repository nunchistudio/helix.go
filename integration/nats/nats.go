package nats

import (
	"strings"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/internal/orchestrator"
	"go.nunchi.studio/helix/service"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

/*
connection represents the nats integration. It respects the integration.Integration
interface and wrap the JetStream interface.
*/
type connection struct {

	// config holds the Config initially passed when creating a new NATS client.
	config *Config

	// nats is the connection made with the NATS server.
	nats *nats.Conn

	// jetstream is the JetStream instance returned by the NATS client, allowing
	// JetStream messaging and stream management.
	jetstream jetstream.JetStream
}

/*
Connect tries to connect to the NATS server given the Config. Returns an error if
Config is not valid or if the connection failed.
*/
func Connect(cfg Config) (JetStream, error) {

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

	// Set the default NATS options.
	opts := []nats.Option{
		nats.Name(orchestrator.Detected.Service()),
		nats.ErrorHandler(asyncErrorHandler),
	}

	// Set TLS options only if enabled in Config.
	if cfg.TLS.Enabled {
		opts = append(opts, nats.ClientCert(cfg.TLS.CertFile, cfg.TLS.KeyFile))

		if len(cfg.TLS.RootCAFiles) > 0 {
			opts = append(opts, nats.RootCAs(cfg.TLS.RootCAFiles...))
		}
	}

	// Try to connect to the NATS servers. Stop here if error validations were
	// encountered.
	addresses := strings.Join(cfg.Addresses, ",")
	conn.nats, err = nats.Connect(addresses, opts...)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})

		return nil, stack
	}

	// Try to return a JetStream instance from NATS. Stop here if error validations
	// were encountered.
	conn.jetstream, err = jetstream.New(conn.nats)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})

		return nil, stack
	}

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, err
	}

	return conn, nil
}
