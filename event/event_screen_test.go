package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventScreenToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Screen
		expected map[string]string
	}{
		{
			input:    Screen{},
			expected: map[string]string{},
		},
		{
			input: Screen{
				Density: 2,
				Width:   12,
				Height:  20,
			},
			expected: map[string]string{
				"event.screen.density": "2",
				"event.screen.width":   "12",
				"event.screen.height":  "20",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventScreenToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventScreenFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.screen.unknown", "anything")
				return m
			},
			expected: &Event{
				Screen: Screen{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.screen.density", "2")
				return m
			},
			expected: &Event{
				Screen: Screen{
					Density: 2,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.screen.width", "12")
				return m
			},
			expected: &Event{
				Screen: Screen{
					Width: 12,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.screen.height", "20")
				return m
			},
			expected: &Event{
				Screen: Screen{
					Height: 20,
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventScreenFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
