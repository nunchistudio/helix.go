package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventReferrerToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Referrer
		expected map[string]string
	}{
		{
			input:    Referrer{},
			expected: map[string]string{},
		},
		{
			input: Referrer{
				Type: "referrer_type_test",
				Name: "referrer_name_test",
				URL:  "referrer_url_test",
				Link: "referrer_link_test",
			},
			expected: map[string]string{
				"event.referrer.type": "referrer_type_test",
				"event.referrer.name": "referrer_name_test",
				"event.referrer.url":  "referrer_url_test",
				"event.referrer.link": "referrer_link_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventReferrerToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventReferrerFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.referrer.unknown", "anything")
				return m
			},
			expected: &Event{
				Referrer: Referrer{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.referrer.type", "referrer_type_test")
				return m
			},
			expected: &Event{
				Referrer: Referrer{
					Type: "referrer_type_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.referrer.name", "referrer_name_test")
				return m
			},
			expected: &Event{
				Referrer: Referrer{
					Name: "referrer_name_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.referrer.url", "referrer_url_test")
				return m
			},
			expected: &Event{
				Referrer: Referrer{
					URL: "referrer_url_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.referrer.link", "referrer_link_test")
				return m
			},
			expected: &Event{
				Referrer: Referrer{
					Link: "referrer_link_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventReferrerFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
