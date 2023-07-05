package temporal

import (
	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

/*
connection represents the temporal integration. It respects the integration.Integration
and Temporal interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new Temporal client.
	config *Config

	// client is the connection made with the Temporal server.
	client client.Client

	// worker holds the Temporal worker, if applicable. It's nil otherwise.
	worker worker.Worker
}

/*
Connect tries to connect to the Temporal server given the Config. Returns an error
if Config is not valid or if the connection failed.
*/
func Connect(cfg Config) (Client, Worker, error) {

	// No need to continue if Config is not valid.
	err := cfg.sanitize()
	if err != nil {
		return nil, nil, err
	}

	// Start to build an error stack, so we can add validations as we go.
	stack := errorstack.New("Failed to initialize integration", errorstack.WithIntegration(identifier))
	conn := &connection{
		config: &cfg,
	}

	// Try to build the tracer.
	t, err := buildTracer(cfg)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})
	}

	// Set the default Temporal config, using cutom logger, context propagator, and
	// tracer.
	var opts = client.Options{
		HostPort:           cfg.Address,
		Namespace:          cfg.Namespace,
		Logger:             new(customlogger),
		ContextPropagators: []workflow.ContextPropagator{new(custompropagator)},
		Interceptors: []interceptor.ClientInterceptor{
			interceptor.NewTracingInterceptor(customtracer{
				Tracer: t,
			}),
		},
	}

	// Set TLS options only if enabled in Config.
	if cfg.TLS.Enabled {
		var validations []errorstack.Validation

		opts.ConnectionOptions.TLS, validations = cfg.TLS.ToStandardTLS()
		if len(validations) > 0 {
			stack.WithValidations(validations...)
		}
	}

	// Try to create the Temporal client.
	conn.client, err = client.Dial(opts)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})
	}

	// Stop here if error validations were encountered. No need to try to create
	// a worker if client is not properly created.
	if stack.HasValidations() {
		return nil, nil, stack
	}

	// Create a Temporal worker if enabled in Config.
	if cfg.Worker.Enabled {
		var optsWorker = worker.Options{
			EnableSessionWorker: true,
		}

		conn.worker = worker.New(conn.client, cfg.Worker.TaskQueue, optsWorker)
		if conn.worker == nil {
			stack.WithValidations(errorstack.Validation{
				Message: "Failed to create worker from client",
			})
		}
	}

	// Stop here if error validations were encountered.
	if stack.HasValidations() {
		return nil, nil, stack
	}

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, nil, err
	}

	// Create the internal Temporal client.
	ic := &iclient{
		config: conn.config,
		client: conn.client,
	}

	// Create the internal Temporal worker.
	iw := &iworker{
		worker: conn.worker,
	}

	return ic, iw, nil
}
