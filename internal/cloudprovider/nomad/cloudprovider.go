package nomad

import (
	"os"

	"go.nunchi.studio/helix/internal/cloudprovider"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
cp is set if the service is running in Nomad, nil otherwise.
*/
var cp cloudprovider.CloudProvider

/*
nomad holds some details about the service currently running in Nomad and
implements the CloudProvider interface.
*/
type nomad struct {
	datacenter string
	jobID      string
	jobName    string
	namespace  string
	region     string
	task       string
}

/*
init populates the cloud provider if the service is running in Nomad.
*/
func init() {
	cp = build()
}

/*
build populates the cloud provider if the service is running in Nomad. Returns
nil otherwise.
*/
func build() cloudprovider.CloudProvider {
	_, exists := os.LookupEnv("NOMAD_JOB_ID")
	if !exists {
		return nil
	}

	n := &nomad{
		datacenter: os.Getenv("NOMAD_DC"),
		jobID:      os.Getenv("NOMAD_JOB_ID"),
		jobName:    os.Getenv("NOMAD_JOB_NAME"),
		namespace:  os.Getenv("NOMAD_NAMESPACE"),
		region:     os.Getenv("NOMAD_REGION"),
		task:       os.Getenv("NOMAD_TASK_NAME"),
	}

	return n
}

/*
Get returns the cloud provider interface for Nomad. Returns nil if not running
in Nomad.
*/
func Get() cloudprovider.CloudProvider {
	return cp
}

/*
String returns the string representation of the Nomad cloud provider.
*/
func (n *nomad) String() string {
	return "nomad"
}

/*
Service returns the service name detected by the cloud provider.
*/
func (n *nomad) Service() string {
	return n.jobName
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
			Key:    "nomad_job_id",
			Type:   zapcore.StringType,
			String: n.jobID,
		},
		zapcore.Field{
			Key:    "nomad_job_name",
			Type:   zapcore.StringType,
			String: n.jobName,
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
		attribute.String("nomad.job_id", n.jobID),
		attribute.String("nomad.job_name", n.jobName),
		attribute.String("nomad.namespace", n.namespace),
		attribute.String("nomad.region", n.region),
		attribute.String("nomad.task", n.task),
		attribute.String("service.name", n.task),
	}

	return attrs
}
