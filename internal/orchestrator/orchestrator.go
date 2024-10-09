package orchestrator

import (
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

/*
Detected is the orchestrator detected by the service package on init. If no known
orchestrator has been detected, this fallbacks to the "unknown" implementation.
*/
var Detected Orchestrator

/*
Orchestrator defines the requirements each orchestrator must meet to be compatible
with the helix.go ecosystem.
*/
type Orchestrator interface {

	// String returns the string representation of the orchestrator.
	//
	// Examples:
	//
	//   "kubernetes"
	//   "nomad"
	//   "unknown"
	String() string

	// Service returns the service name detected by the orchestrator.
	Service() string

	// LoggerFields returns the fields populated by the orchestrator in OpenTelemetry
	// logs.
	LoggerFields() []zap.Field

	// TracerAttributes returns the attributes populated by the orchestrator in
	// OpenTelemetry traces.
	TracerAttributes() []attribute.KeyValue
}
