package bucket

import (
	"fmt"
	"os"

	"go.nunchi.studio/helix/errorstack"

	_ "gocloud.dev/blob/gcsblob"
)

/*
DriverGoogleCloud allows to use Google Cloud Storage as bucket driver. This
driver relies on generic environment variables required by the Google Cloud SDK:

  - GOOGLE_APPLICATION_CREDENTIALS

Config example:

	bucket.Config{
	  Driver: bucket.DriverGoogleCloud,
	  Bucket: "my-bucket",
	}
*/
var DriverGoogleCloud Driver = &driverGoogleCloud{}

/*
driverGoogleCloud is the internal type to use Google Cloud as bucket driver.
*/
type driverGoogleCloud struct{}

/*
string returns the string representation of the Google Cloud bucket driver.
*/
func (d *driverGoogleCloud) string() string {
	return "google"
}

/*
validate ensures Config and environment variables are valid for the Google Cloud
bucket driver.
*/
func (d *driverGoogleCloud) validate(cfg *Config) []errorstack.Validation {
	var validations []errorstack.Validation

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `GOOGLE_APPLICATION_CREDENTIALS` must be set and not be empty",
		})
	}

	return validations
}

/*
url returns the Go Cloud bucket URL of the Google Cloud bucket driver.
*/
func (d *driverGoogleCloud) url(cfg *Config) string {
	path := fmt.Sprintf("gs://%s", cfg.Bucket)

	return path
}
