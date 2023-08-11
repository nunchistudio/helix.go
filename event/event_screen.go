package event

import (
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Screen holds the details about the app's screen from which the event is triggered
from.
*/
type Screen struct {
	Density int64 `json:"density,omitempty"`
	Width   int64 `json:"width,omitempty"`
	Height  int64 `json:"height,omitempty"`
}

/*
injectEventScreenToFlatMap injects values found in a Screen object to a flat map
representation of an Event.
*/
func injectEventScreenToFlatMap(screen Screen, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if screen.Density != 0 {
		flatten["event.screen.density"] = strconv.FormatInt(screen.Density, 10)
	}

	if screen.Width != 0 {
		flatten["event.screen.width"] = strconv.FormatInt(screen.Width, 10)
	}

	if screen.Height != 0 {
		flatten["event.screen.height"] = strconv.FormatInt(screen.Height, 10)
	}
}

/*
applyEventScreenFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Screen. This assumes the Baggage member's
key starts with "event.screen.".
*/
func applyEventScreenFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "density":
		e.Screen.Density, _ = strconv.ParseInt(m.Value(), 10, 0)
	case "width":
		e.Screen.Width, _ = strconv.ParseInt(m.Value(), 10, 0)
	case "height":
		e.Screen.Height, _ = strconv.ParseInt(m.Value(), 10, 0)
	}
}
