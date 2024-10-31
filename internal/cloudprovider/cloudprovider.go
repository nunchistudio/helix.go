package cloudprovider

import (
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

/*
Detected is the cloud provider detected by the service package on init. If no known
cloud provider has been detected, this fallbacks to the "unknown" implementation.
*/
var Detected CloudProvider

/*
CloudProvider defines the requirements each cloud provider must meet to be
compatible with the helix.go ecosystem.
*/
type CloudProvider interface {

	// String returns the string representation of the cloud provider.
	//
	// Examples:
	//
	//   "kubernetes"
	//   "nomad"
	//   "render"
	//   "unknown"
	String() string

	// Service returns the service name detected by the cloud provider.
	Service() string

	// LoggerFields returns the fields populated by the cloud provider in
	// OpenTelemetry logs.
	LoggerFields() []zap.Field

	// TracerAttributes returns the attributes populated by the cloud provider in
	// OpenTelemetry traces.
	TracerAttributes() []attribute.KeyValue
}
