package bucket

import (
	"fmt"
	"os"

	"go.nunchi.studio/helix/errorstack"

	_ "gocloud.dev/blob/s3blob"
)

/*
DriverMinIO allows to use MinIO as bucket driver. This driver relies on generic
environment variables required by the MinIO SDK:

  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY
  - AWS_ENDPOINT

Config example:

	bucket.Config{
	  Driver: bucket.DriverMinIO,
	  Bucket: "my-bucket",
	}
*/
var DriverMinIO Driver = &driverMinIO{}

/*
driverMinIO is the internal type to use MinIO as bucket driver.
*/
type driverMinIO struct{}

/*
string returns the string representation of the MinIO bucket driver.
*/
func (d *driverMinIO) string() string {
	return "minio"
}

/*
validate ensures Config and environment variables are valid for the MinIO bucket
driver.
*/
func (d *driverMinIO) validate(cfg *Config) []errorstack.Validation {
	var validations []errorstack.Validation

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `AWS_ACCESS_KEY_ID` must be set and not be empty",
		})
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `AWS_SECRET_ACCESS_KEY` must be set and not be empty",
		})
	}

	if os.Getenv("AWS_ENDPOINT") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `AWS_ENDPOINT` must be set and not be empty",
		})
	}

	return validations
}

/*
url returns the Go Cloud bucket URL of the MinIO bucket driver.
*/
func (d *driverMinIO) url(cfg *Config) string {
	path := fmt.Sprintf("s3://%s?endpoint=%s&disableSSL=true&s3ForcePathStyle=true", cfg.Bucket, os.Getenv("AWS_ENDPOINT"))

	return path
}
