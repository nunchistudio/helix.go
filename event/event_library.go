package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Library holds the details of the SDK used by the client executing the event.
*/
type Library struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

/*
injectEventLibraryToFlatMap injects values found in a Library object to a flat
map representation of an Event.
*/
func injectEventLibraryToFlatMap(library Library, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if library.Name != "" {
		flatten["event.library.name"] = library.Name
	}

	if library.Version != "" {
		flatten["event.library.version"] = library.Version
	}
}

/*
applyEventLibraryFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Library. This assumes the Baggage member's
key starts with "event.library.".
*/
func applyEventLibraryFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "name":
		e.Library.Name = m.Value()
	case "version":
		e.Library.Version = m.Value()
	}
}
