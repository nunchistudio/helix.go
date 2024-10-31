package render

import (
	"os"

	"go.nunchi.studio/helix/internal/cloudprovider"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
cp is set if the service is running in Render, nil otherwise.
*/
var cp cloudprovider.CloudProvider

/*
render holds some details about the service currently running in Render and
implements the CloudProvider interface.
*/
type render struct {
	instanceID  string
	serviceID   string
	serviceName string
	serviceType string
}

/*
init populates the cloud provider if the service is running in Render.
*/
func init() {
	cp = build()
}

/*
build populates the cloud provider if the service is running in Render. Returns
nil otherwise.
*/
func build() cloudprovider.CloudProvider {
	_, exists := os.LookupEnv("RENDER")
	if !exists {
		return nil
	}

	n := &render{
		instanceID:  os.Getenv("RENDER_INSTANCE_ID"),
		serviceID:   os.Getenv("RENDER_SERVICE_ID"),
		serviceName: os.Getenv("RENDER_SERVICE_NAME"),
		serviceType: os.Getenv("RENDER_SERVICE_TYPE"),
	}

	return n
}

/*
Get returns the cloud provider interface for Render. Returns nil if not running
in Render.
*/
func Get() cloudprovider.CloudProvider {
	return cp
}

/*
String returns the string representation of the Render cloud provider.
*/
func (q *render) String() string {
	return "render"
}

/*
Service returns the service name detected by the cloud provider.
*/
func (q *render) Service() string {
	return q.serviceID
}

/*
LoggerFields returns OpenTelemetry fields for logs when running in Render.
*/
func (q *render) LoggerFields() []zap.Field {
	fields := []zap.Field{
		zapcore.Field{
			Key:    "render_instance_id",
			Type:   zapcore.StringType,
			String: q.instanceID,
		},
		zapcore.Field{
			Key:    "render_service_id",
			Type:   zapcore.StringType,
			String: q.serviceID,
		},
		zapcore.Field{
			Key:    "render_service_name",
			Type:   zapcore.StringType,
			String: q.serviceName,
		},
		zapcore.Field{
			Key:    "render_service_type",
			Type:   zapcore.StringType,
			String: q.serviceType,
		},
	}

	return fields
}

/*
TracerAttributes returns OpenTelemetry attributes for traces when running in
Render.
*/
func (q *render) TracerAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("render.instance_id", q.instanceID),
		attribute.String("render.service_id", q.serviceID),
		attribute.String("render.service_name", q.serviceName),
		attribute.String("render.service_type", q.serviceType),
		attribute.String("service.name", q.serviceID),
	}

	return attrs
}
