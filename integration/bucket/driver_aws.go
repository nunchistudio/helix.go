package bucket

import (
	"fmt"
	"os"

	"go.nunchi.studio/helix/errorstack"

	_ "gocloud.dev/blob/s3blob"
)

/*
DriverAWS allows to use AWS S3 as bucket driver. This driver relies on generic
environment variables required by the AWS SDK (v2):

  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY
  - AWS_REGION

Config example:

	bucket.Config{
	  Driver: bucket.DriverAWS,
	  Bucket: "my-bucket",
	}
*/
var DriverAWS Driver = &driverAWS{}

/*
driverAWS is the internal type to use AWS as bucket driver.
*/
type driverAWS struct{}

/*
string returns the string representation of the AWS bucket driver.
*/
func (d *driverAWS) string() string {
	return "aws"
}

/*
validate ensures Config and environment variables are valid for the AWS S3 bucket
driver.
*/
func (d *driverAWS) validate(cfg *Config) []errorstack.Validation {
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

	if os.Getenv("AWS_REGION") == "" {
		validations = append(validations, errorstack.Validation{
			Message: "Environment variable `AWS_REGION` must be set and not be empty",
		})
	}

	return validations
}

/*
url returns the Go Cloud bucket URL of the AWS S3 bucket driver.
*/
func (d *driverAWS) url(cfg *Config) string {
	path := fmt.Sprintf("s3://%s?region=%s&awssdk=v2", cfg.Bucket, os.Getenv("AWS_REGION"))

	return path
}
