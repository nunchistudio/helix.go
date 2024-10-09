package nomad

import (
	"os"

	"go.nunchi.studio/helix/internal/orchestrator"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
orch is set if the service is running in Nomad, nil otherwise.
*/
var orch orchestrator.Orchestrator

/*
nomad holds some details about the service currently running in Nomad and implements
the Orchestrator interface.
*/
type nomad struct {
	datacenter string
	job        string
	namespace  string
	region     string
	task       string
}

/*
init populates the orchestrator if the service is running in Nomad.
*/
func init() {
	orch = build()
}

/*
build populates the orchestrator if the service is running in Nomad. Returns
nil otherwise.
*/
func build() orchestrator.Orchestrator {
	_, exists := os.LookupEnv("NOMAD_JOB_ID")
	if !exists {
		return nil
	}

	n := &nomad{
		datacenter: os.Getenv("NOMAD_DC"),
		job:        os.Getenv("NOMAD_JOB_NAME"),
		namespace:  os.Getenv("NOMAD_NAMESPACE"),
		region:     os.Getenv("NOMAD_REGION"),
		task:       os.Getenv("NOMAD_TASK_NAME"),
	}

	return n
}

/*
Get returns the orchestrator interface for Nomad. Returns nil if not running
in Nomad.
*/
func Get() orchestrator.Orchestrator {
	return orch
}

/*
String returns the string representation of the Nomad orchestrator.
*/
func (n *nomad) String() string {
	return "nomad"
}

/*
Service returns the service name detected by the orchestrator.
*/
func (n *nomad) Service() string {
	return n.task
}

/*
LoggerFields returns OpenTelemetry fields for logs when running in Nomad.
*/
func (n *nomad) LoggerFields() []zap.Field {
	fields := []zap.Field{
		zapcore.Field{
			Key:    "nomad_datacenter",
			Type:   zapcore.StringType,
			String: n.datacenter,
		},
		zapcore.Field{
			Key:    "nomad_job",
			Type:   zapcore.StringType,
			String: n.job,
		},
		zapcore.Field{
			Key:    "nomad_namespace",
			Type:   zapcore.StringType,
			String: n.namespace,
		},
		zapcore.Field{
			Key:    "nomad_region",
			Type:   zapcore.StringType,
			String: n.region,
		},
		zapcore.Field{
			Key:    "nomad_task",
			Type:   zapcore.StringType,
			String: n.task,
		},
	}

	return fields
}

/*
TracerAttributes returns OpenTelemetry attributes for traces when running in
Nomad.
*/
func (n *nomad) TracerAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("nomad.datacenter", n.datacenter),
		attribute.String("nomad.job", n.job),
		attribute.String("nomad.namespace", n.namespace),
		attribute.String("nomad.region", n.region),
		attribute.String("nomad.task", n.task),
		attribute.String("service.name", n.task),
	}

	return attrs
}
