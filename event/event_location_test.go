package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventLocationToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Location
		expected map[string]string
	}{
		{
			input:    Location{},
			expected: map[string]string{},
		},
		{
			input: Location{
				City:      "location_city_test",
				Country:   "location_country_test",
				Region:    "location_region_test",
				Latitude:  45.916,
				Longitude: 6.133,
				Speed:     50,
			},
			expected: map[string]string{
				"event.location.city":      "location_city_test",
				"event.location.country":   "location_country_test",
				"event.location.region":    "location_region_test",
				"event.location.latitude":  "45.916000",
				"event.location.longitude": "6.133000",
				"event.location.speed":     "50.000000",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventLocationToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventLocationFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.unknown", "anything")
				return m
			},
			expected: &Event{
				Location: Location{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.city", "location_city_test")
				return m
			},
			expected: &Event{
				Location: Location{
					City: "location_city_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.country", "location_country_test")
				return m
			},
			expected: &Event{
				Location: Location{
					Country: "location_country_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.region", "location_region_test")
				return m
			},
			expected: &Event{
				Location: Location{
					Region: "location_region_test",
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.latitude", "45.916000")
				return m
			},
			expected: &Event{
				Location: Location{
					Latitude: 45.916,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.longitude", "6.133000")
				return m
			},
			expected: &Event{
				Location: Location{
					Longitude: 6.133,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.location.speed", "50.000000")
				return m
			},
			expected: &Event{
				Location: Location{
					Speed: 50,
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventLocationFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
