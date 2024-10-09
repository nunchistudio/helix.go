package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventAppToFlatMap(t *testing.T) {
	testcases := []struct {
		input    App
		expected map[string]string
	}{
		{
			input:    App{},
			expected: map[string]string{},
		},
		{
			input: App{
				Name:    "app_name_test",
				Version: "app_version_test",
				BuildID: "app_buildid_test",
			},
			expected: map[string]string{
				"event.app.name":     "app_name_test",
				"event.app.version":  "app_version_test",
				"event.app.build_id": "app_buildid_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventAppToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventAppFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.app.unknown", "anything")
				return m
			},
			expected: &Event{
				App: App{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.app.name", "app_name_test")
				return m
			},
			expected: &Event{
				App: App{
					Name: "app_name_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.app.version", "app_version_test")
				return m
			},
			expected: &Event{
				App: App{
					Version: "app_version_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.app.build_id", "app_buildid_test")
				return m
			},
			expected: &Event{
				App: App{
					BuildID: "app_buildid_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventAppFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
