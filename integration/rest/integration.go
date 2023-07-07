package rest

import (
	"context"
	"net/http"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"
	"go.nunchi.studio/helix/internal/orchestrator"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

/*
Ensure *rest complies to the integration.Integration type.
*/
var _ integration.Integration = (*rest)(nil)

/*
String returns the string representation of the HTTP REST integration.
*/
func (r *rest) String() string {
	return identifier
}

/*
Start starts the HTTP server of the HTTP REST integration.
*/
func (r *rest) Start(ctx context.Context) error {
	stack := errorstack.New("Failed to start HTTP server", errorstack.WithIntegration(identifier))

	// Wrap the built-in HTTP handler with the one given by the user, if applicable.
	var h http.Handler = r.bun
	if r.config.Middleware != nil {
		h = r.config.Middleware(r.bun)
	}

	// Wrap the handler previously built with the one designed for OpenTelemetry
	// traces.
	h = otelhttp.NewHandler(h, orchestrator.Detected.Service(),
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	// Create the HTTP server with the given configuration and the handler built.
	r.server = &http.Server{
		Addr:    r.config.Address,
		Handler: h,
	}

	// Start the HTTP server with or without TLS depending on the Config, and catch
	// unexpected errors.
	var err error
	if r.config.TLS.Enabled {
		err = r.server.ListenAndServeTLS(r.config.TLS.CertFile, r.config.TLS.KeyFile)
	} else {
		err = r.server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return stack
	}

	return nil
}

/*
Close tries to gracefully close the HTTP server.
*/
func (r *rest) Close(ctx context.Context) error {
	stack := errorstack.New("Failed to gracefully close HTTP server", errorstack.WithIntegration(identifier))

	err := r.server.Shutdown(ctx)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return stack
	}

	return nil
}

/*
Status always returns a `200` status.
*/
func (r *rest) Status(ctx context.Context) (int, error) {
	return 200, nil
}
