package bucket

import (
	"fmt"

	"go.nunchi.studio/helix/errorstack"

	_ "gocloud.dev/blob/fileblob"
)

/*
DriverLocal allows to use local files as bucket driver.

Config example:

	bucket.Config{
	  Driver: bucket.DriverLocal,
	  Bucket: "/path/to/folder",
	}
*/
var DriverLocal Driver = &driverLocal{}

/*
driverLocal is the internal type to use local files as bucket driver.
*/
type driverLocal struct{}

/*
string returns the string representation of the local bucket driver.
*/
func (d *driverLocal) string() string {
	return "local"
}

/*
validate ensures Config and environment variables are valid for the local bucket
driver.
*/
func (d *driverLocal) validate(cfg *Config) []errorstack.Validation {
	var validations []errorstack.Validation

	return validations
}

/*
url returns the Go Cloud bucket URL of the local bucket driver.
*/
func (d *driverLocal) url(cfg *Config) string {
	path := fmt.Sprintf("file://%s", cfg.Bucket)

	return path
}
