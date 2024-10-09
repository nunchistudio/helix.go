package bucket

import (
	"go.nunchi.studio/helix/errorstack"
)

/*
Driver allows to set the driver to use for the bucket integration.
*/
type Driver interface {

	// string returns the string representation of the driver.
	string() string

	// validate ensures Config and environment variables are valid for the driver.
	validate(cfg *Config) []errorstack.Validation

	// url returns the Go Cloud bucket URL of the bucket driver.
	url(cfg *Config) string
}
