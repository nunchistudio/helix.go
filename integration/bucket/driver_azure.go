package bucket

import (
	"fmt"
	"os"

	"go.nunchi.studio/helix/errorstack"

	_ "gocloud.dev/blob/azureblob"
)

/*
DriverAzure allows to use Azure Blob Storage as bucket driver. This driver relies
on generic environment variables required by the Azure SDK:

  - AZURE_STORAGE_ACCOUNT
  - AZURE_STORAGE_KEY || AZURE_STORAGE_SAS_TOKEN

Config example:

	bucket.Config{
	  Driver: bucket.DriverAzure,
	  Bucket: "my-container",
	}
*/
var DriverAzure Driver = &driverAzure{}

/*
driverAzure is the internal type to use Azure as bucket driver.
*/
type driverAzure struct{}

/*
string returns the string representation of the Azure bucket driver.
*/
func (d *driverAzure) string() string {
	return "azure"
}

/*
validate ensures Config and environment variables are valid for the Azure bucket
driver.
*/
func (d *driverAzure) validate(cfg *Config) []errorstack.Validation {
	var validations []errorstack.Validation

	if os.Getenv("AZURE_STORAGE_ACCOUNT") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `AZURE_STORAGE_ACCOUNT` must be set and not be empty",
		})
	}

	if os.Getenv("AZURE_STORAGE_KEY") == "" && os.Getenv("AZURE_STORAGE_SAS_TOKEN") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "One of environment variable `AZURE_STORAGE_KEY` or `AZURE_STORAGE_SAS_TOKEN` must be set and not be empty",
		})
	}

	return validations
}

/*
url returns the Go Cloud bucket URL of the Azure bucket driver.
*/
func (d *driverAzure) url(cfg *Config) string {
	path := fmt.Sprintf("azblob://%s", cfg.Bucket)

	return path
}
