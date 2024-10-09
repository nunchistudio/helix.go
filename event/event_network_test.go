package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestInjectEventNetworkToFlatMap(t *testing.T) {
	testcases := []struct {
		input    Network
		expected map[string]string
	}{
		{
			input:    Network{},
			expected: map[string]string{},
		},
		{
			input: Network{
				Bluetooth: true,
				Cellular:  false,
				WIFI:      true,
				Carrier:   "network_carrier_test",
			},
			expected: map[string]string{
				"event.network.bluetooth": "true",
				"event.network.wifi":      "true",
				"event.network.carrier":   "network_carrier_test",
			},
		},
	}

	for _, tc := range testcases {
		var actual = make(map[string]string)
		for k, v := range tc.expected {
			actual[k] = v
		}

		injectEventNetworkToFlatMap(tc.input, tc.expected)

		assert.Equal(t, tc.expected, actual)
	}
}

func TestApplyEventNetworkFromBaggageMember(t *testing.T) {
	testcases := []struct {
		input    func() baggage.Member
		expected *Event
	}{
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.network.unknown", "anything")
				return m
			},
			expected: &Event{
				Network: Network{},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.network.bluetooth", "true")
				return m
			},
			expected: &Event{
				Network: Network{
					Bluetooth: true,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.network.cellular", "true")
				return m
			},
			expected: &Event{
				Network: Network{
					Cellular: true,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.network.wifi", "true")
				return m
			},
			expected: &Event{
				Network: Network{
					WIFI: true,
				},
			},
		},
		{
			input: func() baggage.Member {
				m, _ := baggage.NewMember("event.network.carrier", "network_carrier_test")
				return m
			},
			expected: &Event{
				Network: Network{
					Carrier: "network_carrier_test",
				},
			},
		},
	}

	for _, tc := range testcases {
		actual := new(Event)
		applyEventNetworkFromBaggageMember(tc.input(), actual)

		assert.Equal(t, tc.expected, actual)
	}
}
