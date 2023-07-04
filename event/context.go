package event

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.nunchi.studio/helix/internal/contextkey"

	"go.opentelemetry.io/otel/baggage"
)

/*
EventFromContext returns the Event found in the context passed, if any. If no
Event has been found, it tries to find and build one if a Baggage was found in
the context. Returns true if an Event has been found, false otherwise.
*/
func EventFromContext(ctx context.Context) (Event, bool) {
	var e Event

	e, ok := ctx.Value(contextkey.Event).(Event)
	if !ok {
		return eventFromBaggage(baggage.FromContext(ctx))
	}

	return e, true
}

/*
ContextWithEvent returns a copy of the context passed with the Event associated
to it.
*/
func ContextWithEvent(ctx context.Context, e Event) context.Context {
	return context.WithValue(ctx, contextkey.Event, e)
}

/*
eventFromBaggage returns the Event found in the Baggage passed, if any. Returns
true if an Event has been found, false otherwise.

Example:

	map[string]string = {
	  "event.name"                     = "subscribed",
	  "event.user_id"                  = "user_2N6YZQLcYy2SPtmHiII69yHp0WE,
	  "event.params.filters.0"         = "a",
	  "event.params.filters.1"         = "b",
	  "event.params.filters.2"         = "c",
	  "event.subscription.id"          = "sub_2N6YZQXgQAv87zMmvlHxePCSsRs",
	  "event.subscription.customer_id" = "cus_2N6YZMi3sBDPQBZrZJoYBwhNQNv",
	  "event.subscription.plan_id"     = "plan_2N6YZSE1SkWT9DrlXlswLhJ5K5Q",
	}

Will produce:

	Event{
	  Name:   "subscribed",
	  UserID: "user_2N6YZQLcYy2SPtmHiII69yHp0WE,
	  Params: url.Values{
	    "filters": []string{"a", "b", "c"},
		},
	  Subscription: Subscription{
	    ID:         "sub_2N6YZQXgQAv87zMmvlHxePCSsRs",
	    CustomerID: "cus_2N6YZMi3sBDPQBZrZJoYBwhNQNv",
	    PlanID:     "plan_2N6YZSE1SkWT9DrlXlswLhJ5K5Q",
	  },
	}
*/
func eventFromBaggage(b baggage.Baggage) (Event, bool) {
	var e = Event{}

	e.Name = b.Member("event.name").Value()
	if e.Name == "" {
		return e, false
	}

	for _, m := range b.Members() {
		if !strings.HasPrefix(m.Key(), "event.") {
			continue
		}

		if strings.HasPrefix(m.Key(), "event.meta.") {
			if e.Meta == nil {
				e.Meta = make(map[string]string)
			}

			e.Meta[strings.TrimPrefix(m.Key(), "event.meta.")] = m.Value()
			continue
		}

		if strings.HasPrefix(m.Key(), "event.params.") {
			if e.Params == nil {
				e.Params = make(url.Values)
			}

			key := strings.TrimPrefix(m.Key(), "event.params.")
			split := strings.Split(key, ".")

			e.Params.Add(split[0], m.Value())
			continue
		}

		if strings.HasPrefix(m.Key(), "event.subscription.flags.") {
			if e.Subscription.Flags == nil {
				e.Subscription.Flags = make(map[string]string)
			}

			e.Subscription.Flags[strings.TrimPrefix(m.Key(), "event.subscription.flags.")] = m.Value()
			continue
		}

		switch strings.TrimPrefix(m.Key(), "event.") {
		case "is_anonymous":
			e.IsAnonymous, _ = strconv.ParseBool(m.Value())
		case "user_id":
			e.UserID = m.Value()
		case "group_id":
			e.GroupID = m.Value()

		case "subscription.id":
			e.Subscription.ID = m.Value()
		case "subscription.customer_id":
			e.Subscription.CustomerID = m.Value()
		case "subscription.plan_id":
			e.Subscription.PlanID = m.Value()
		case "subscription.usage":
			e.Subscription.Usage = m.Value()
		case "subscription.increment_by":
			e.Subscription.IncrementBy, _ = strconv.ParseFloat(m.Value(), 64)

		case "app.name":
			e.App.Name = m.Value()
		case "app.version":
			e.App.Version = m.Value()
		case "app.build_id":
			e.App.BuildID = m.Value()

		case "library.name":
			e.Library.Name = m.Value()
		case "library.version":
			e.Library.Version = m.Value()

		case "campaign.name":
			e.Campaign.Name = m.Value()
		case "campaign.source":
			e.Campaign.Source = m.Value()
		case "campaign.medium":
			e.Campaign.Medium = m.Value()
		case "campaign.term":
			e.Campaign.Term = m.Value()
		case "campaign.content":
			e.Campaign.Content = m.Value()

		case "referrer.type":
			e.Referrer.Type = m.Value()
		case "referrer.name":
			e.Referrer.Name = m.Value()
		case "referrer.url":
			e.Referrer.URL = m.Value()
		case "referrer.link":
			e.Referrer.Link = m.Value()

		case "cloud.provider":
			e.Cloud.Provider = m.Value()
		case "cloud.service":
			e.Cloud.Service = m.Value()
		case "cloud.region":
			e.Cloud.Region = m.Value()
		case "cloud.project_id":
			e.Cloud.ProjectID = m.Value()
		case "cloud.account_id":
			e.Cloud.AccountID = m.Value()

		case "device.id":
			e.Device.ID = m.Value()
		case "device.manufacturer":
			e.Device.Manufacturer = m.Value()
		case "device.model":
			e.Device.Model = m.Value()
		case "device.name":
			e.Device.Name = m.Value()
		case "device.type":
			e.Device.Type = m.Value()
		case "device.version":
			e.Device.Version = m.Value()
		case "device.advertising_id":
			e.Device.AdvertisingID = m.Value()

		case "os.name":
			e.OS.Name = m.Value()
		case "os.arch":
			e.OS.Arch = m.Value()
		case "os.version":
			e.OS.Version = m.Value()

		case "location.city":
			e.Location.City = m.Value()
		case "location.country":
			e.Location.Country = m.Value()
		case "location.region":
			e.Location.Region = m.Value()
		case "location.latitude":
			e.Location.Latitude, _ = strconv.ParseFloat(m.Value(), 64)
		case "location.longitude":
			e.Location.Longitude, _ = strconv.ParseFloat(m.Value(), 64)
		case "location.speed":
			e.Location.Speed, _ = strconv.ParseFloat(m.Value(), 64)

		case "network.bluetooth":
			e.Network.Bluetooth, _ = strconv.ParseBool(m.Value())
		case "network.cellular":
			e.Network.Cellular, _ = strconv.ParseBool(m.Value())
		case "network.wifi":
			e.Network.WIFI, _ = strconv.ParseBool(m.Value())
		case "network.carrier":
			e.Network.Carrier = m.Value()

		case "page.path":
			e.Page.Path = m.Value()
		case "page.referrer":
			e.Page.Referrer = m.Value()
		case "page.search":
			e.Page.Search = m.Value()
		case "page.title":
			e.Page.Title = m.Value()
		case "page.url":
			e.Page.URL = m.Value()

		case "screen.density":
			e.Screen.Density, _ = strconv.ParseInt(m.Value(), 10, 0)
		case "screen.width":
			e.Screen.Width, _ = strconv.ParseInt(m.Value(), 10, 0)
		case "screen.height":
			e.Screen.Height, _ = strconv.ParseInt(m.Value(), 10, 0)

		case "ip":
			e.IP = net.ParseIP(m.Value())
		case "locale":
			e.Locale = m.Value()
		case "timezone":
			e.Timezone = m.Value()
		case "user_agent":
			e.UserAgent = m.Value()
		case "timestamp":
			e.Timestamp, _ = time.Parse(time.RFC3339Nano, m.Value())
		}
	}

	return e, true
}
