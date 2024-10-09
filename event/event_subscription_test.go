package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventSubscriptionsToFlatMap(t *testing.T) {
	testcases := []struct {
		input    []Subscription
		expected map[string]string
	}{
		{
			input:    []Subscription{},
			expected: map[string]string{},
		},
		{
			input: []Subscription{
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
				{
					ID:          "subscription_1_id_test",
					CustomerID:  "subscription_1_customerid_test",
					PlanID:      "subscription_1_planid_test",
					Usage:       "subscription_1_usage_test",
					IncrementBy: 1.25,
					Flags: map[string]string{
						"version": "b",
					},
				},
			},
			expected: map[string]string{
				"event.subscriptions[0].id":            "subscription_0_id_test",
				"event.subscriptions[0].customer_id":   "subscription_0_customerid_test",
				"event.subscriptions[0].plan_id":       "subscription_0_planid_test",
				"event.subscriptions[0].usage":         "subscription_0_usage_test",
				"event.subscriptions[0].increment_by":  "1.000000",
				"event.subscriptions[0].flags.version": "a",
				"event.subscriptions[1].id":            "subscription_1_id_test",
				"event.subscriptions[1].customer_id":   "subscription_1_customerid_test",
				"event.subscriptions[1].plan_id":       "subscription_1_planid_test",
				"event.subscriptions[1].usage":         "subscription_1_usage_test",
				"event.subscriptions[1].increment_by":  "1.250000",
				"event.subscriptions[1].flags.version": "b",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventSubscriptionsToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventSubscriptionsFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.unknown", "anything")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.id", "subscription_0_id_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						ID: "subscription_0_id_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.customer_id", "subscription_0_customerid_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						CustomerID: "subscription_0_customerid_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.plan_id", "subscription_0_planid_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						PlanID: "subscription_0_planid_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.usage", "subscription_0_usage_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						Usage: "subscription_0_usage_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.increment_by", "1.000000")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						IncrementBy: 1,
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.0.flags.version", "a")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{
						Flags: map[string]string{
							"version": "a",
						},
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.unknown", "anything")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.id", "subscription_1_id_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						ID: "subscription_1_id_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.customer_id", "subscription_1_customerid_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						CustomerID: "subscription_1_customerid_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.plan_id", "subscription_1_planid_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						PlanID: "subscription_1_planid_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.usage", "subscription_1_usage_test")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						Usage: "subscription_1_usage_test",
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.increment_by", "1.250000")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						IncrementBy: 1.25,
					},
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.subscriptions.1.flags.version", "b")
				return m
			},
			expected: &Event{
				Subscriptions: []Subscription{
					{},
					{
						Flags: map[string]string{
							"version": "b",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventSubscriptionsFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
