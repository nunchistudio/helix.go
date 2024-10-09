package event

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

/*
Location holds the details about the user's location.
*/
type Location struct {
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

/*
injectEventLocationToFlatMap injects values found in a Location object to a flat
map representation of an Event.
*/
func injectEventLocationToFlatMap(location Location, flatten map[string]string) {
	if flatten == nil {
		flatten = make(map[string]string)
	}

	if location.City != "" {
		flatten["event.location.city"] = location.City
	}

	if location.Country != "" {
		flatten["event.location.country"] = location.Country
	}

	if location.Region != "" {
		flatten["event.location.region"] = location.Region
	}

	if location.Latitude != 0 {
		flatten["event.location.latitude"] = fmt.Sprintf("%f", location.Latitude)
	}

	if location.Longitude != 0 {
		flatten["event.location.longitude"] = fmt.Sprintf("%f", location.Longitude)
	}

	if location.Speed != 0 {
		flatten["event.location.speed"] = fmt.Sprintf("%f", location.Speed)
	}
}

/*
applyEventLocationFromBaggageMember extracts the value of a Baggage member given
its key and applies it to an Event's Location. This assumes the Baggage member's
key starts with "event.location.".
*/
func applyEventLocationFromBaggageMember(m baggage.Member, e *Event) {
	split := strings.Split(m.Key(), ".")

	switch split[2] {
	case "city":
		e.Location.City = m.Value()
	case "country":
		e.Location.Country = m.Value()
	case "region":
		e.Location.Region = m.Value()
	case "latitude":
		e.Location.Latitude, _ = strconv.ParseFloat(m.Value(), 64)
	case "longitude":
		e.Location.Longitude, _ = strconv.ParseFloat(m.Value(), 64)
	case "speed":
		e.Location.Speed, _ = strconv.ParseFloat(m.Value(), 64)
	}
}
