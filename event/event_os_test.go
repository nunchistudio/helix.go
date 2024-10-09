package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventOSToFlatMap(t *testing.T) {
	testcases := []struct {
		input    OS
		expected map[string]string
	}{
		{
			input:    OS{},
			expected: map[string]string{},
		},
		{
			input: OS{
				Name:    "app_name_test",
				Arch:    "app_arch_test",
				Version: "app_version_test",
			},
			expected: map[string]string{
				"event.os.name":    "app_name_test",
				"event.os.arch":    "app_arch_test",
				"event.os.version": "app_version_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventOSToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventOSFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.os.unknown", "anything")
				return m
			},
			expected: &Event{
				OS: OS{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.os.name", "app_name_test")
				return m
			},
			expected: &Event{
				OS: OS{
					Name: "app_name_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.os.arch", "app_arch_test")
				return m
			},
			expected: &Event{
				OS: OS{
					Arch: "app_arch_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.os.version", "app_version_test")
				return m
			},
			expected: &Event{
				OS: OS{
					Version: "app_version_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventOSFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
