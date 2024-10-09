package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Cloud holds the details about the cloud provider from which the client is executing
the event.
*/
type Cloud struct {
	Provider  string `json:"provider,omitempty"`
	Service   string `json:"service,omitempty"`
	Region    string `json:"region,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

/*
injectEventCloudToFlatMap injects values found in a Cloud object to a flat map
representation of an Event.
*/
func injectEventCloudToFlatMap(cloud Cloud, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if cloud.Provider != "" {
		flatten["event.cloud.provider"] = cloud.Provider
	}

	if cloud.Service != "" {
		flatten["event.cloud.service"] = cloud.Service
	}

	if cloud.Region != "" {
		flatten["event.cloud.region"] = cloud.Region
	}

	if cloud.ProjectID != "" {
		flatten["event.cloud.project_id"] = cloud.ProjectID
	}

	if cloud.AccountID != "" {
		flatten["event.cloud.account_id"] = cloud.AccountID
	}
}

/*
applyEventCloudFromBaggageMember extracts the value of a Baggage member given its
key and applies it to an Event's Cloud. This assumes the Baggage member's key
starts with "event.cloud.".
*/
func applyEventCloudFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "provider":
		e.Cloud.Provider = m.Value()
	case "service":
		e.Cloud.Service = m.Value()
	case "region":
		e.Cloud.Region = m.Value()
	case "project_id":
		e.Cloud.ProjectID = m.Value()
	case "account_id":
		e.Cloud.AccountID = m.Value()
	}
}
