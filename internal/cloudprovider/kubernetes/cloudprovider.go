package kubernetes

import (
	"os"

	"go.nunchi.studio/helix/internal/cloudprovider"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
cp is set if the service is running in Kubernetes, nil otherwise.
*/
var cp cloudprovider.CloudProvider

/*
kubernetes holds some details about the service currently running in Kubernetes
and implements the CloudProvider interface.
*/
type kubernetes struct {
	namespace string
	pod       string
}

/*
init populates the cloud provider if the service is running in Kubernetes.
*/
func init() {
	cp = build()
}

/*
build populates the cloud provider if the service is running in Kubernetes.
Returns nil otherwise.
*/
func build() cloudprovider.CloudProvider {
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
Get returns the cloud provider interface for Kubernetes. Returns nil if not
running in Kubernetes.
*/
func Get() cloudprovider.CloudProvider {
	return cp
}

/*
String returns the string representation of the Kubernetes cloud provider.
*/
func (k *kubernetes) String() string {
	return "kubernetes"
}

/*
Service returns the service name detected by the cloud provider.
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
