package clickhouse

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

/*
ClickHouse exposes an opinionated way to interact with ClickHouse, by bringing
automatic distributed tracing as well as error recording within traces.
*/
type ClickHouse interface {
	Query(ctx context.Context, query string, args ...any) (driver.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) driver.Row
	PrepareBatch(ctx context.Context, query string) (Batch, error)
	Exec(ctx context.Context, query string, args ...any) error
	AsyncInsert(ctx context.Context, query string, wait bool) error
}

/*
connection represents the clickhouse integration. It respects the
integration.Integration and ClickHouse interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new ClickHouse client.
	config *Config

	// client is the connection made with the ClickHouse server.
	client driver.Conn
}

/*
Connect tries to connect to the ClickHouse server given the Config. Returns an
error if Config is not valid or if the connection failed.
*/
func Connect(cfg Config) (ClickHouse, error) {

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

	// Set the default ClickHouse options.
	var opts = &clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     cfg.Addresses,
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	}

	// Set TLS options only if enabled in Config.
	if cfg.TLS.Enabled {
		var validations []errorstack.Validation

		opts.TLS, validations = cfg.TLS.ToStandardTLS()
		if len(validations) > 0 {
			stack.WithValidations(validations...)
		}
	}

	// Try to connect to the ClickHouse servers.
	conn.client, err = clickhouse.Open(opts)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})
	}

	// Stop here if error validations were encountered.
	if stack.HasValidations() {
		return nil, stack
	}

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

/*
Query executes a query returning multiple rows.

It automatically handles tracing and error recording.
*/
func (conn *connection) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: QueryRows", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to query rows", err)
		}
	}()

	rows, err := conn.client.Query(ctx, query, args...)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return rows, err
}

/*
QueryRow executes a query returning a single row.

It automatically handles tracing.
*/
func (conn *connection) QueryRow(ctx context.Context, query string, args ...any) driver.Row {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: QueryRow", humanized))
	defer span.End()

	row := conn.client.QueryRow(ctx, query, args...)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return row
}

/*
PrepareBatch prepares and bind a new batch.
*/
func (conn *connection) PrepareBatch(ctx context.Context, query string) (Batch, error) {
	pb, err := conn.client.PrepareBatch(ctx, query)
	if err != nil {
		return nil, err
	}

	b := &batch{
		config: conn.config,
		client: pb,
	}

	return b, nil
}

/*
Exec executes a query and does not expect any rows to be returned.

It automatically handles tracing and error recording.
*/
func (conn *connection) Exec(ctx context.Context, query string, args ...any) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Exec", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to execute query", err)
		}
	}()

	err = conn.client.Exec(ctx, query, args...)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return err
}

/*
AsyncInsert asynchronously inserts a query.

It automatically handles tracing and error recording.
*/
func (conn *connection) AsyncInsert(ctx context.Context, query string, wait bool) error {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: AsyncInsert", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to insert asynchronously", err)
		}
	}()

	err = conn.client.AsyncInsert(ctx, query, wait)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)
	span.SetBoolAttribute(fmt.Sprintf("%s.async_insert.wait", identifier), wait)

	return err
}
