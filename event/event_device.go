package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Device holds the details about the user's device.
*/
type Device struct {
	ID            string `json:"id,omitempty"`
	Manufacturer  string `json:"manufacturer,omitempty"`
	Model         string `json:"model,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Version       string `json:"version,omitempty"`
	AdvertisingID string `json:"advertising_id,omitempty"`
}

/*
injectEventDeviceToFlatMap injects values found in a Device object to a flat map
representation of an Event.
*/
func injectEventDeviceToFlatMap(device Device, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if device.ID != "" {
		flatten["event.device.id"] = device.ID
	}

	if device.Manufacturer != "" {
		flatten["event.device.manufacturer"] = device.Manufacturer
	}

	if device.Model != "" {
		flatten["event.device.model"] = device.Model
	}

	if device.Name != "" {
		flatten["event.device.name"] = device.Name
	}

	if device.Type != "" {
		flatten["event.device.type"] = device.Type
	}

	if device.Version != "" {
		flatten["event.device.version"] = device.Version
	}

	if device.AdvertisingID != "" {
		flatten["event.device.advertising_id"] = device.AdvertisingID
	}
}

/*
applyEventDeviceFromBaggageMember extracts the value of a Baggage member given its
key and applies it to an Event's Device. This assumes the Baggage member's key
starts with "event.device.".
*/
func applyEventDeviceFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "id":
		e.Device.ID = m.Value()
	case "manufacturer":
		e.Device.Manufacturer = m.Value()
	case "model":
		e.Device.Model = m.Value()
	case "name":
		e.Device.Name = m.Value()
	case "type":
		e.Device.Type = m.Value()
	case "version":
		e.Device.Version = m.Value()
	case "advertising_id":
		e.Device.AdvertisingID = m.Value()
	}
}
