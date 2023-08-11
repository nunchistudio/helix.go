package event

import (
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Network holds the details about the user's network.
*/
type Network struct {
	Bluetooth bool   `json:"bluetooth,omitempty"`
	Cellular  bool   `json:"cellular,omitempty"`
	WIFI      bool   `json:"wifi,omitempty"`
	Carrier   string `json:"carrier,omitempty"`
}

/*
injectEventNetworkToFlatMap injects values found in a Network object to a flat
map representation of an Event.
*/
func injectEventNetworkToFlatMap(network Network, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if network.Bluetooth {
		flatten["event.network.bluetooth"] = strconv.FormatBool(network.Bluetooth)
	}

	if network.Cellular {
		flatten["event.network.cellular"] = strconv.FormatBool(network.Cellular)
	}

	if network.WIFI {
		flatten["event.network.wifi"] = strconv.FormatBool(network.WIFI)
	}

	if network.Carrier != "" {
		flatten["event.network.carrier"] = network.Carrier
	}
}

/*
applyEventNetworkFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Network. This assumes the Baggage member's
key starts with "event.network.".
*/
func applyEventNetworkFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "bluetooth":
		e.Network.Bluetooth, _ = strconv.ParseBool(m.Value())
	case "cellular":
		e.Network.Cellular, _ = strconv.ParseBool(m.Value())
	case "wifi":
		e.Network.WIFI, _ = strconv.ParseBool(m.Value())
	case "carrier":
		e.Network.Carrier = m.Value()
	}
}
