package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventCloudToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Cloud
		expected map[string]string
	}{
		{
			input:    Cloud{},
			expected: map[string]string{},
		},
		{
			input: Cloud{
				Provider:  "cloud_provider_test",
				Service:   "cloud_service_test",
				Region:    "cloud_region_test",
				ProjectID: "cloud_projectid_test",
				AccountID: "cloud_accountid_test",
			},
			expected: map[string]string{
				"event.cloud.provider":   "cloud_provider_test",
				"event.cloud.service":    "cloud_service_test",
				"event.cloud.region":     "cloud_region_test",
				"event.cloud.project_id": "cloud_projectid_test",
				"event.cloud.account_id": "cloud_accountid_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventCloudToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventCloudFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.unknown", "anything")
				return m
			},
			expected: &Event{
				Cloud: Cloud{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.provider", "cloud_provider_test")
				return m
			},
			expected: &Event{
				Cloud: Cloud{
					Provider: "cloud_provider_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.service", "cloud_service_test")
				return m
			},
			expected: &Event{
				Cloud: Cloud{
					Service: "cloud_service_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.region", "cloud_region_test")
				return m
			},
			expected: &Event{
				Cloud: Cloud{
					Region: "cloud_region_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.project_id", "cloud_projectid_test")
				return m
			},
			expected: &Event{
				Cloud: Cloud{
					ProjectID: "cloud_projectid_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.cloud.account_id", "cloud_accountid_test")
				return m
			},
			expected: &Event{
				Cloud: Cloud{
					AccountID: "cloud_accountid_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventCloudFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
