package kubernetes

import (
	"os"

	"go.nunchi.studio/helix/internal/orchestrator"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
orch is set if the service is running in Kubernetes, nil otherwise.
*/
var orch orchestrator.Orchestrator

/*
kubernetes holds some details about the service currently running in Kubernetes
and implements the Orchestrator interface.
*/
type kubernetes struct {
	namespace string
	pod       string
}

/*
init populates the orchestrator if the service is running in Kubernetes.
*/
func init() {
	orch = build()
}

/*
build populates the orchestrator if the service is running in Kubernetes. Returns
nil otherwise.
*/
func build() orchestrator.Orchestrator {
	_, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !exists {
		return nil
	}

	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil
	}

	k := &kubernetes{
		namespace: string(ns),
		pod:       os.Getenv("HOSTNAME"),
	}

	return k
}

/*
Get returns the orchestrator interface for Kubernetes. Returns nil if not running
in Kubernetes.
*/
func Get() orchestrator.Orchestrator {
	return orch
}

/*
String returns the string representation of the Kubernetes orchestrator.
*/
func (k *kubernetes) String() string {
	return "kubernetes"
}

/*
Service returns the service name detected by the orchestrator.
*/
func (k *kubernetes) Service() string {
	return k.pod
}

/*
LoggerFields returns OpenTelemetry fields for logs when running in Kubernetes.
*/
func (k *kubernetes) LoggerFields() []zap.Field {
	fields := []zap.Field{
		zapcore.Field{
			Key:    "kubernetes_namespace",
			Type:   zapcore.StringType,
			String: k.namespace,
		},
		zapcore.Field{
			Key:    "kubernetes_pod",
			Type:   zapcore.StringType,
			String: k.pod,
		},
	}

	return fields
}

/*
TracerAttributes returns OpenTelemetry attributes for traces when running in
Kubernetes.
*/
func (k *kubernetes) TracerAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("kubernetes.namespace", k.namespace),
		attribute.String("kubernetes.pod", k.pod),
		attribute.String("service.name", k.pod),
	}

	return attrs
}
