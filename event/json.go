package event

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
Key is the key that shall be present in a JSON-encoded value representing
an Event.

Example:

	{
	  "key": "value",
	  "event": {
	    "name": "subscribed"
	  }
	}
*/
const Key string = "event"

/*
EventFromJSON returns the Event found at the "event" key in the JSON-encoded data
passed, if any. Returns true if an Event has been found, false otherwise.
*/
func EventFromJSON(input json.RawMessage) (Event, bool) {
	var e Event
	var mapped map[string]any
	err := json.Unmarshal(input, &mapped)
	if err != nil {
		return e, false
	}

	v, exists := mapped[Key]
	if !exists {
		return e, false
	}

	b, err := json.Marshal(v)
	if err != nil {
		return e, false
	}

	err = json.Unmarshal(b, &e)
	if err != nil {
		return e, false
	}

	return e, true
}

/*
ToFlatMap returns a flatten map for a given Event. Keys are prefixed with "event.",
and struct level are seperated by a ".". All values are stringified.

This is primarly designed for the telemetry packages, allowing to pass contextual
information about an event using Go's context or HTTP headers, but can be useful
in some other use cases.

Example:

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

Will produce:

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
*/
func ToFlatMap(e Event) map[string]string {
	var flatten = make(map[string]string)

	flatten["event.name"] = e.Name

	if e.Meta != nil {
		for k, v := range e.Meta {
			flatten[fmt.Sprintf("event.meta.%s", k)] = v
		}
	}

	if e.Params != nil {
		for k, v := range e.Params {
			split := strings.Split(k, ".")
			for i, s := range v {
				flatten[fmt.Sprintf("event.params.%s.%d", split[0], i)] = s
			}
		}
	}

	flatten["event.is_anonymous"] = strconv.FormatBool(e.IsAnonymous)
	flatten["event.user_id"] = e.UserID
	flatten["event.group_id"] = e.GroupID

	flatten["event.subscription.id"] = e.Subscription.ID
	flatten["event.subscription.customer_id"] = e.Subscription.CustomerID
	flatten["event.subscription.plan_id"] = e.Subscription.PlanID
	flatten["event.subscription.usage"] = e.Subscription.Usage
	flatten["event.subscription.increment_by"] = fmt.Sprintf("%f", e.Subscription.IncrementBy)
	if e.Subscription.Flags != nil {
		for k, v := range e.Subscription.Flags {
			flatten[fmt.Sprintf("event.subscription.flags.%s", k)] = v
		}
	}

	flatten["event.app.name"] = e.App.Name
	flatten["event.app.version"] = e.App.Version
	flatten["event.app.build_id"] = e.App.BuildID

	flatten["event.library.name"] = e.Library.Name
	flatten["event.library.version"] = e.Library.Version

	flatten["event.campaign.name"] = e.Campaign.Name
	flatten["event.campaign.source"] = e.Campaign.Source
	flatten["event.campaign.medium"] = e.Campaign.Medium
	flatten["event.campaign.term"] = e.Campaign.Term
	flatten["event.campaign.content"] = e.Campaign.Content

	flatten["event.referrer.type"] = e.Referrer.Type
	flatten["event.referrer.name"] = e.Referrer.Name
	flatten["event.referrer.url"] = e.Referrer.URL
	flatten["event.referrer.link"] = e.Referrer.Link

	flatten["event.cloud.provider"] = e.Cloud.Provider
	flatten["event.cloud.service"] = e.Cloud.Service
	flatten["event.cloud.region"] = e.Cloud.Region
	flatten["event.cloud.project_id"] = e.Cloud.ProjectID
	flatten["event.cloud.account_id"] = e.Cloud.AccountID

	flatten["event.device.id"] = e.Device.ID
	flatten["event.device.manufacturer"] = e.Device.Manufacturer
	flatten["event.device.model"] = e.Device.Model
	flatten["event.device.name"] = e.Device.Name
	flatten["event.device.type"] = e.Device.Type
	flatten["event.device.version"] = e.Device.Version
	flatten["event.device.advertising_id"] = e.Device.AdvertisingID

	flatten["event.os.name"] = e.OS.Name
	flatten["event.os.arch"] = e.OS.Arch
	flatten["event.os.version"] = e.OS.Version

	flatten["event.location.city"] = e.Location.City
	flatten["event.location.country"] = e.Location.Country
	flatten["event.location.region"] = e.Location.Region
	flatten["event.location.latitude"] = fmt.Sprintf("%f", e.Location.Latitude)
	flatten["event.location.longitude"] = fmt.Sprintf("%f", e.Location.Longitude)
	flatten["event.location.speed"] = fmt.Sprintf("%f", e.Location.Speed)

	flatten["event.network.bluetooth"] = strconv.FormatBool(e.Network.Bluetooth)
	flatten["event.network.cellular"] = strconv.FormatBool(e.Network.Cellular)
	flatten["event.network.wifi"] = strconv.FormatBool(e.Network.WIFI)
	flatten["event.network.carrier"] = e.Network.Carrier

	flatten["event.page.path"] = e.Page.Path
	flatten["event.page.referrer"] = e.Page.Referrer
	flatten["event.page.search"] = e.Page.Search
	flatten["event.page.title"] = e.Page.Title
	flatten["event.page.url"] = e.Page.URL

	flatten["event.screen.density"] = strconv.FormatInt(e.Screen.Density, 10)
	flatten["event.screen.width"] = strconv.FormatInt(e.Screen.Width, 10)
	flatten["event.screen.height"] = strconv.FormatInt(e.Screen.Height, 10)

	flatten["event.ip"] = e.IP.String()
	flatten["event.locale"] = e.Locale
	flatten["event.timezone"] = e.Timezone
	flatten["event.user_agent"] = e.UserAgent

	if !e.Timestamp.IsZero() {
		flatten["event.timestamp"] = e.Timestamp.Format(time.RFC3339Nano)
	}

	for k, v := range flatten {
		if v == "" || v == "false" || v == "0" || v == "0E+00" || v == "0.000000" || v == "<nil>" {
			delete(flatten, k)
		}
	}

	return flatten
}
