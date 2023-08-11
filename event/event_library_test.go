package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventLibraryToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Library
		expected map[string]string
	}{
		{
			input:    Library{},
			expected: map[string]string{},
		},
		{
			input: Library{
				Name:    "library_name_test",
				Version: "library_version_test",
			},
			expected: map[string]string{
				"event.library.name":    "library_name_test",
				"event.library.version": "library_version_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventLibraryToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventLibraryFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.library.unknown", "anything")
				return m
			},
			expected: &Event{
				Library: Library{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.library.name", "library_name_test")
				return m
			},
			expected: &Event{
				Library: Library{
					Name: "library_name_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.library.version", "library_version_test")
				return m
			},
			expected: &Event{
				Library: Library{
					Version: "library_version_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventLibraryFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
