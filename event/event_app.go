package event

import (
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
App holds the details about the client application executing the event.
*/
type App struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	BuildID string `json:"build_id,omitempty"`
}

/*
injectEventAppToFlatMap injects values found in an App object to a flat map
representation of an Event.
*/
func injectEventAppToFlatMap(app App, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if app.Name != "" {
		flatten["event.app.name"] = app.Name
	}

	if app.Version != "" {
		flatten["event.app.version"] = app.Version
	}

	if app.BuildID != "" {
		flatten["event.app.build_id"] = app.BuildID
	}
}

/*
applyEventAppFromBaggageMember extracts the value of a Baggage member given its
key and applies it to an Event's App. This assumes the Baggage member's key starts
with "event.app.".
*/
func applyEventAppFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "name":
		e.App.Name = m.Value()
	case "version":
		e.App.Version = m.Value()
	case "build_id":
		e.App.BuildID = m.Value()
	}
}
