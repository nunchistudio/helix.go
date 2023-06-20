package event

/*
Subscription holds the details about the account/customer from which the event
has been triggered. It's useful for tracking customer usages.
*/
type Subscription struct {
	ID          string            `json:"id,omitempty"`
	CustomerID  string            `json:"customer_id,omitempty"`
	PlanID      string            `json:"plan_id,omitempty"`
	Usage       string            `json:"usage,omitempty"`
	IncrementBy float64           `json:"increment_by,omitempty"`
	Flags       map[string]string `json:"flags,omitempty"`
}

/*
App holds the details about the client application executing the event.
*/
type App struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	BuildID string `json:"build_id,omitempty"`
}

/*
Library holds the details of the SDK used by the client executing the event.
*/
type Library struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

/*
Campaign holds the details about the marketing campaign from which a client is
executing the event from.
*/
type Campaign struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Term    string `json:"term,omitempty"`
	Content string `json:"content,omitempty"`
}

/*
Referrer holds the details about the marketing referrer from which a client is
executing the event from.
*/
type Referrer struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Link string `json:"link,omitempty"`
}

/*
Cloud holds the details about the cloud provider from which the client is executing
the event.
*/
type Cloud struct {
	Provider  string `json:"provider,omitempty"`
	Service   string `json:"service,omitempty"`
	Region    string `json:"region,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

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
OS holds the details about the user's OS.
*/
type OS struct {
	Name    string `json:"name,omitempty"`
	Arch    string `json:"arch,omitempty"`
	Version string `json:"version,omitempty"`
}

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
Network holds the details about the user's network.
*/
type Network struct {
	Bluetooth bool   `json:"bluetooth,omitempty"`
	Cellular  bool   `json:"cellular,omitempty"`
	WIFI      bool   `json:"wifi,omitempty"`
	Carrier   string `json:"carrier,omitempty"`
}

/*
Page holds the details about the webpage from which the event is triggered from.
*/
type Page struct {
	Path     string `json:"path,omitempty"`
	Referrer string `json:"referrer,omitempty"`
	Search   string `json:"search,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
}

/*
Screen holds the details about the app's screen from which the event is triggered
from.
*/
type Screen struct {
	Density int64 `json:"density,omitempty"`
	Width   int64 `json:"width,omitempty"`
	Height  int64 `json:"height,omitempty"`
}
