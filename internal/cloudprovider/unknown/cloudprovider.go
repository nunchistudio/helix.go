package unknown

import (
	"os"
	"path/filepath"

	"go.nunchi.studio/helix/internal/cloudprovider"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
cp is always set since the "unknown" cloud provider is used as fallback in case
no other cloud provider has been detected.
*/
var cp cloudprovider.CloudProvider

/*
unknown holds some details about the service currently running and implements the
CloudProvider interface.
*/
type unknown struct {
	name string
}

/*
init populates the cloud provider as a fallback cloud provider.
*/
func init() {
	cp = build()
}

/*
build populates the cloudprovider.Detected as a fallback cloud provider. If no
cloud provider is returned it means an internal error occurred while finding the
path to the Go executable currently being run, fallbacks to a static string if
necessary. This should never happen.
*/
func build() cloudprovider.CloudProvider {
	var name string = "helix"

	path, err := os.Executable()
	if err == nil {
		name = filepath.Base(path)
	}

	u := &unknown{
		name: name,
	}

	return u
}

/*
Get returns the fallback cloud provider interface.
*/
func Get() cloudprovider.CloudProvider {
	return cp
}

/*
String returns the string representation of the unknown cloud provider.
*/
func (u *unknown) String() string {
	return "unknown"
}

/*
Service returns the service name detected by the cloud provider.
*/
func (u *unknown) Service() string {
	return u.name
}

/*
LoggerFields returns basic OpenTelemetry fields for logs.
*/
func (u *unknown) LoggerFields() []zap.Field {
	fields := []zap.Field{
		zapcore.Field{
			Key:    "service_name",
			Type:   zapcore.StringType,
			String: u.name,
		},
	}

	return fields
}

/*
TracerAttributes returns basic OpenTelemetry attributes for traces.
*/
func (u *unknown) TracerAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("service.name", u.name),
	}

	return attrs
}
