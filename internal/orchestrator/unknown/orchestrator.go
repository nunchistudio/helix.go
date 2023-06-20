package unknown

import (
	"os"
	"path/filepath"

	"go.nunchi.studio/helix/internal/orchestrator"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
orch is always set since the "unknown" orchestrator is used as fallback in case
no other orchestrator has been detected.
*/
var orch orchestrator.Orchestrator

/*
unknown holds some details about the service currently running and implements the
Orchestrator interface.
*/
type unknown struct {
	name string
}

/*
init populates the orchestrator as a fallback orchestrator.
*/
func init() {
	orch = build()
}

/*
build populates the orchestrator as a fallback orchestrator. If no orchestrator
is returned it means an internal error occured while finding the path to the
Go executable currently being run, fallbacks to a static string if necessary.
This should never happen.
*/
func build() orchestrator.Orchestrator {
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
Get returns the fallback orchestrator interface.
*/
func Get() orchestrator.Orchestrator {
	return orch
}

/*
String returns the string representation of the unknown orchestrator.
*/
func (u *unknown) String() string {
	return "unknown"
}

/*
Service returns the service name detected by the orchestrator.
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
