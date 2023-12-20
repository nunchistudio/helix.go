package bucket

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"
	"go.nunchi.studio/helix/telemetry/trace"

	"gocloud.dev/blob"
)

/*
Bucket exposes an opinionated and standardized way to interact with buckets across
different providers through drivers.
*/
type Bucket interface {
	Exists(ctx context.Context, key string) bool
	Read(ctx context.Context, key string) ([]byte, error)
	Write(ctx context.Context, key string, value []byte, opts *OptionsWrite) error
	Copy(ctx context.Context, srcKey string, dstKey string) error
	Delete(ctx context.Context, key string) error
}

/*
connection represents the bucket integration. It respects the integration.Integration
and Bucket interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new Bucket client.
	config *Config

	// client is the connection made with the Bucket client.
	client *blob.Bucket
}

/*
Connect tries to create a Bucket client given the Config. Returns an error if
Config is not valid or if the initialization failed.
*/
func Connect(cfg Config) (Bucket, error) {

	// No need to continue if Config is not valid.
	err := cfg.sanitize()
	if err != nil {
		return nil, err
	}

	// Start to build an error stack, so we can add validations as we go.
	stack := errorstack.New("Failed to initialize integration", errorstack.WithIntegration(identifier))
	conn := &connection{
		config: &cfg,
	}

	// Try to create the Bucket connection, using the URL returned by the driver.
	conn.client, err = blob.OpenBucket(context.Background(), cfg.Driver.url(&cfg))
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return nil, stack
	}

	// Use a prefixed bucket, if applicable.
	if cfg.Subfolder != "" {
		conn.client = blob.PrefixedBucket(conn.client, cfg.Subfolder)
	}

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

/*
Exists returns true if a blob exists at key, false otherwise.

It automatically handles tracing.
*/
func (conn *connection) Exists(ctx context.Context, key string) bool {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Exists", humanized))
	defer span.End()

	exists, _ := conn.client.Exists(ctx, key)
	setDefaultAttributes(span, conn.config)
	setKeyAttributes(span, key)

	return exists
}

/*
Read reads the blob at key and returns its byte representation.

It automatically handles tracing and error recording.
*/
func (conn *connection) Read(ctx context.Context, key string) ([]byte, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Read", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to read blob", err)
		}
	}()

	value, err := conn.client.ReadAll(ctx, key)
	setDefaultAttributes(span, conn.config)
	setKeyAttributes(span, key)

	return value, err
}

/*
Write writes bytes represenation of blob at key, with some optional options.

It automatically handles tracing and error recording.
*/
func (conn *connection) Write(ctx context.Context, key string, value []byte, opts *OptionsWrite) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Write", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to write blob", err)
		}
	}()

	var write *blob.WriterOptions
	if opts != nil {
		write = &blob.WriterOptions{
			CacheControl:       opts.CacheControl,
			ContentDisposition: opts.ContentDisposition,
			ContentEncoding:    opts.ContentEncoding,
			ContentLanguage:    opts.ContentLanguage,
			ContentType:        opts.ContentType,
			ContentMD5:         opts.ContentMD5,
			Metadata:           opts.Metadata,
		}
	}

	err = conn.client.WriteAll(ctx, key, value, write)
	setDefaultAttributes(span, conn.config)
	setKeyAttributes(span, key)

	return err
}

/*
Copy copies blob from srcKey to dstKey.

It automatically handles tracing and error recording.
*/
func (conn *connection) Copy(ctx context.Context, srcKey string, dstKey string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Copy", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to copy blob", err)
		}
	}()

	err = conn.client.Copy(ctx, dstKey, srcKey, nil)
	setDefaultAttributes(span, conn.config)
	span.SetStringAttribute(fmt.Sprintf("%s.key_source", identifier), srcKey)
	span.SetStringAttribute(fmt.Sprintf("%s.key_destination", identifier), dstKey)

	return err
}

/*
Delete deletes existing blob at key.

It automatically handles tracing and error recording.
*/
func (conn *connection) Delete(ctx context.Context, key string) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Delete", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to delete blob", err)
		}
	}()

	err = conn.client.Delete(ctx, key)
	setDefaultAttributes(span, conn.config)
	setKeyAttributes(span, key)

	return err
}
