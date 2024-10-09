package event

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Event
		expected map[string]string
	}{
		{
			input:    Event{},
			expected: map[string]string{},
		},
		{
			input: Event{
				Name: "name_test",
				Meta: map[string]string{
					"hello": "world",
					"this":  "works",
				},
				Params: url.Values{
					"query": []string{"a", "b"},
				},
				IsAnonymous: false,
				Subscriptions: []Subscription{
					{
						ID:          "subscription_0_id_test",
						CustomerID:  "subscription_0_customerid_test",
						PlanID:      "subscription_0_planid_test",
						Usage:       "subscription_0_usage_test",
						IncrementBy: 1,
						Flags: map[string]string{
							"version": "a",
						},
					},
				},
			},
			expected: map[string]string{
				"event.name":                           "name_test",
				"event.meta.hello":                     "world",
				"event.meta.this":                      "works",
				"event.params.query[0]":                "a",
				"event.params.query[1]":                "b",
				"event.subscriptions[0].id":            "subscription_0_id_test",
				"event.subscriptions[0].customer_id":   "subscription_0_customerid_test",
				"event.subscriptions[0].plan_id":       "subscription_0_planid_test",
				"event.subscriptions[0].usage":         "subscription_0_usage_test",
				"event.subscriptions[0].increment_by":  "1.000000",
				"event.subscriptions[0].flags.version": "a",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestExtractEventFromBaggage(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Baggage
		expected Event
	}{
		{
			input: func() baggage.Baggage {
				name, _ := baggage.NewMember("event.name", "name_test")
				metaHelloWorld, _ := baggage.NewMember("event.meta.hello", "world")
				metaThisWorks, _ := baggage.NewMember("event.meta.this", "works")

				b, _ := baggage.New(name, metaHelloWorld, metaThisWorks)
				return b
			},
			expected: Event{
				Name: "name_test",
				Meta: map[string]string{
					"hello": "world",
					"this":  "works",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := extractEventFromBaggage(tc.input())

		assert.Equal(t, tc.expected, actual)
	}
}
