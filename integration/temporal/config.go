package temporal

import (
	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/integration"

	"go.temporal.io/sdk/converter"
)

/*
Config is used to configure the Temporal integration.
*/
type Config struct {

	// Address is the Temporal server address to connect to.
	//
	// Default:
	//
	//   "127.0.0.1:7233"
	Address string `json:"address"`

	// Namespace sets the namespace to connect to.
	//
	// Default:
	//
	//   "default"
	Namespace string `json:"namespace"`

	// DataConverter customizes serialization/deserialization of arguments in
	// Temporal.
	DataConverter converter.DataConverter `json:"-"`

	// Worker configures a Temporal worker if the helix service should run as worker
	// for Temporal.
	Worker ConfigWorker `json:"worker"`

	// TLSConfig configures TLS to communicate with the Temporal server.
	TLS integration.ConfigTLS `json:"tls"`
}

/*
ConfigWorker configures a Temporal worker for the helix service running. When
enabled, this starts a Temporal worker for the given task queue and namespace
(set in Config).
*/
type ConfigWorker struct {

	// Enabled creates a Temporal worker, to run workflows and activities.
	Enabled bool `json:"enabled"`

	// TaskQueue is the task queue name you use to identify your client worker,
	// also identifies group of workflow and activity implementations that are hosted
	// by a single worker process.
	//
	// Required when enabled.
	TaskQueue string `json:"taskqueue,omitempty"`

	// WorkerActivitiesPerSecond sets the rate limiting on number of activities that
	// can be executed per second per worker. This can be used to limit resources
	// used by the worker.
	//
	// Notice that the number is represented in float, so that you can set it to
	// less than 1 if needed. For example, set the number to 0.1 means you want
	// your activity to be executed once for every 10 seconds. This can be used to
	// protect down stream services from flooding.
	//
	// Default:
	//
	//   100 000
	WorkerActivitiesPerSecond float64 `json:"worker_activities_per_second,omitempty"`

	// TaskQueueActivitiesPerSecond sets the rate limiting on number of activities
	// that can be executed per second. This is managed by the server and controls
	// activities per second for your entire taskqueue.
	//
	// Notice that the number is represented in float, so that you can set it to
	// less than 1 if needed. For example, set the number to 0.1 means you want
	// your activity to be executed once for every 10 seconds. This can be used to
	// protect down stream services from flooding.
	//
	// Default:
	//
	//   100 000
	TaskQueueActivitiesPerSecond float64 `json:"taskqueue_activities_per_second,omitempty"`
}

/*
sanitize sets default values - when applicable - and validates the configuration.
Returns an error if configuration is not valid.
*/
func (cfg *Config) sanitize() error {
	stack := errorstack.New("Failed to validate configuration", errorstack.WithIntegration(identifier))

	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:7233"
	}

	if cfg.Namespace == "" {
		cfg.Namespace = "default"
	}

	if cfg.Worker.Enabled {
		if cfg.Worker.TaskQueue == "" {
			stack.WithValidations(errorstack.Validation{
				Message: "TaskQueue must be set and not be empty",
				Path:    []string{"Config", "Worker", "TaskQueue"},
			})
		}
	}

	stack.WithValidations(cfg.TLS.Sanitize()...)
	if stack.HasValidations() {
		return stack
	}

	return nil
}
