package qovery

import (
	"os"

	"go.nunchi.studio/helix/internal/cloudprovider"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
cp is set if the service is running in Qovery, nil otherwise.
*/
var cp cloudprovider.CloudProvider

/*
qovery holds some details about the service currently running in Qovery and
implements the CloudProvider interface.
*/
type qovery struct {
	region              string
	projectID           string
	environmentID       string
	applicationID       string
	deploymentID        string
	kubernetesCluster   string
	kubernetesNamespace string
	kubernetesPod       string
}

/*
init populates the cloud provider if the service is running in Qovery.
*/
func init() {
	cp = build()
}

/*
build populates the cloud provider if the service is running in Qovery. Returns
nil otherwise.
*/
func build() cloudprovider.CloudProvider {
	_, exists := os.LookupEnv("QOVERY_APPLICATION_ID")
	if !exists {
		return nil
	}

	n := &qovery{
		region:              os.Getenv("QOVERY_CLOUD_PROVIDER_REGION"),
		projectID:           os.Getenv("QOVERY_PROJECT_ID"),
		environmentID:       os.Getenv("QOVERY_ENVIRONMENT_ID"),
		applicationID:       os.Getenv("QOVERY_APPLICATION_ID"),
		deploymentID:        os.Getenv("QOVERY_DEPLOYMENT_ID"),
		kubernetesCluster:   os.Getenv("QOVERY_KUBERNETES_CLUSTER_NAME"),
		kubernetesNamespace: os.Getenv("QOVERY_KUBERNETES_NAMESPACE_NAME"),
		kubernetesPod:       os.Getenv("HOSTNAME"),
	}

	return n
}

/*
Get returns the cloud provider interface for Qovery. Returns nil if not running
in Qovery.
*/
func Get() cloudprovider.CloudProvider {
	return cp
}

/*
String returns the string representation of the Qovery cloud provider.
*/
func (q *qovery) String() string {
	return "qovery"
}

/*
Service returns the service name detected by the cloud provider.
*/
func (q *qovery) Service() string {
	return q.kubernetesPod
}

/*
LoggerFields returns OpenTelemetry fields for logs when running in Qovery.
*/
func (q *qovery) LoggerFields() []zap.Field {
	fields := []zap.Field{
		zapcore.Field{
			Key:    "qovery_region",
			Type:   zapcore.StringType,
			String: q.region,
		},
		zapcore.Field{
			Key:    "qovery_project_id",
			Type:   zapcore.StringType,
			String: q.projectID,
		},
		zapcore.Field{
			Key:    "qovery_environment_id",
			Type:   zapcore.StringType,
			String: q.environmentID,
		},
		zapcore.Field{
			Key:    "qovery_application_id",
			Type:   zapcore.StringType,
			String: q.applicationID,
		},
		zapcore.Field{
			Key:    "qovery_deployment_id",
			Type:   zapcore.StringType,
			String: q.deploymentID,
		},
		zapcore.Field{
			Key:    "kubernetes_cluster",
			Type:   zapcore.StringType,
			String: q.kubernetesCluster,
		},
		zapcore.Field{
			Key:    "kubernetes_namespace",
			Type:   zapcore.StringType,
			String: q.kubernetesNamespace,
		},
		zapcore.Field{
			Key:    "kubernetes_pod",
			Type:   zapcore.StringType,
			String: q.kubernetesPod,
		},
	}

	return fields
}

/*
TracerAttributes returns OpenTelemetry attributes for traces when running in
Qovery.
*/
func (q *qovery) TracerAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("qovery.region", q.region),
		attribute.String("qovery.project_id", q.projectID),
		attribute.String("qovery.environment_id", q.environmentID),
		attribute.String("qovery.application_id", q.applicationID),
		attribute.String("qovery.deployment_id", q.deploymentID),
		attribute.String("kubernetes.cluster", q.kubernetesCluster),
		attribute.String("kubernetes.namespace", q.kubernetesNamespace),
		attribute.String("kubernetes.pod", q.kubernetesPod),
		attribute.String("service.name", q.kubernetesPod),
	}

	return attrs
}
