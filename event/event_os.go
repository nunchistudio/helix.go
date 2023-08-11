package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
OS holds the details about the user's OS.
*/
type OS struct {
	Name    string `json:"name,omitempty"`
	Arch    string `json:"arch,omitempty"`
	Version string `json:"version,omitempty"`
}

/*
injectEventOSToFlatMap injects values found in an OS object to a flat map
representation of an Event.
*/
func injectEventOSToFlatMap(os OS, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if os.Name != "" {
		flatten["event.os.name"] = os.Name
	}

	if os.Arch != "" {
		flatten["event.os.arch"] = os.Arch
	}

	if os.Version != "" {
		flatten["event.os.version"] = os.Version
	}
}

/*
applyEventOSFromBaggageMember extracts the value of a Baggage member given its
key and applies it to an Event's OS. This assumes the Baggage member's key starts
with "event.os.".
*/
func applyEventOSFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "name":
		e.OS.Name = m.Value()
	case "arch":
		e.OS.Arch = m.Value()
	case "version":
		e.OS.Version = m.Value()
	}
}
