package postgres

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/service"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*
PostgreSQL exposes an opinionated way to interact with PostgreSQL, by bringing
automatic distributed tracing as well as error recording within traces.
*/
type PostgreSQL interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Prepare(ctx context.Context, id string, query string) (*pgconn.StatementDescription, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults
	WaitForNotification(ctx context.Context) (*pgconn.Notification, error)
}

/*
connection represents the postgres integration. It respects the
integration.Integration and PostgreSQL interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new PostgreSQL client.
	config *Config

	// client is the connection made with the PostgreSQL server.
	client *pgx.Conn
}

/*
Connect tries to connect to the PostgreSQL server given the Config. Returns an
error if Config is not valid or if the connection failed.
*/
func Connect(cfg Config) (PostgreSQL, error) {

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

	// Set the default PostgreSQL options.
	address := fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Address, cfg.Database)
	opts, err := pgx.ParseConfig(address)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: normalizeErrorMessage(err),
		})

		return nil, stack
	}

	// Wrap and apply the notification. We wrap so end-users don't have access to
	// the underlying PostgreSQL connection, but also so one day we could potentially
	// add logic such as tracing.
	if cfg.OnNotification != nil {
		opts.OnNotification = func(pc *pgconn.PgConn, notif *pgconn.Notification) {
			cfg.OnNotification(notif)
		}
	}

	// Set TLS options only if enabled in Config.
	if cfg.TLS.Enabled {
		opts.Config.TLSConfig, stack.Validations = cfg.TLS.ToStandardTLS()
	}

	// Try to connect to the PostgreSQL servers.
	conn.client, err = pgx.ConnectConfig(context.Background(), opts)
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
BeginTx starts a transaction. Unlike database/sql, the context only affects the
begin command: there is no auto-rollback on context cancellation.

It automatically handles tracing and error recording.
*/
func (conn *connection) BeginTx(ctx context.Context, opts pgx.TxOptions) (Tx, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Begin", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to begin transaction", err)
		}
	}()

	client, err := conn.client.BeginTx(ctx, opts)
	setDefaultAttributes(span, conn.config)

	tx := &transaction{
		config: conn.config,
		client: client,
	}

	return tx, err
}

/*
Exec executes a SQL query. It can be either a prepared statement name or an SQL
string. Arguments should be referenced positionally from the query string as
$1, $2, etc.

It automatically handles tracing and error recording.
*/
func (conn *connection) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Exec", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to execute query", err)
		}
	}()

	stmt, err := conn.client.Exec(ctx, query, args...)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return stmt, err
}

/*
Prepare creates a prepared statement with a unique name. The query can contain
placeholders or bound parameters. These placeholders are referenced positional
as $1, $2, etc.

Prepare is idempotent: it is safe to call Prepare multiple times with the same
name and query arguments. This allows a code path to Prepare and Query/Exec
without concern for if the statement has already been prepared.

It automatically handles tracing and error recording.
*/
func (conn *connection) Prepare(ctx context.Context, id string, query string) (*pgconn.StatementDescription, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Prepare", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to prepare statement", err)
		}
	}()

	stmt, err := conn.client.Prepare(ctx, id, query)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return stmt, err
}

/*
Query sends a query to the server and returns a Rows to read the results. Only
errors encountered sending the query and initializing Rows will be returned.

Err() on the returned Rows must be checked after the Rows is closed to determine
if the query executed successfully.

The returned Rows must be closed before the connection can be used again. It is
safe to attempt to read from the returned Rows even if an error is returned. The
error will be the available in rows.Err() after rows are closed. It is allowed to
ignore the error returned from Query and handle it in Rows.

It is possible for a query to return one or more rows before encountering an
error. In most cases the rows should be collected before processing rather than
processed while receiving each row. This avoids the possibility of the application
processing rows from a query that the server rejected. The pgx.CollectRows function
is useful here.

An implementor of QueryRewriter may be passed as the first element of args. It
can rewrite the query and change or replace args. For example, pgx.NamedArgs is
QueryRewriter that implements named arguments.

It automatically handles tracing and error recording.
*/
func (conn *connection) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
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
QueryRow is a convenience wrapper over Query. Any error that occurs while querying
is deferred until calling Scan on the returned Row. That Row will error with
pgx.ErrNoRows if no rows are returned.

It automatically handles tracing.
*/
func (conn *connection) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: QueryRow", humanized))
	defer span.End()

	row := conn.client.QueryRow(ctx, query, args...)
	setDefaultAttributes(span, conn.config)
	setQueryAttributes(span, query)

	return row
}

/*
SendBatch sends all queued queries to the server at once. All queries are run in
an implicit transaction unless explicit transaction control statements are executed.
The returned BatchResults must be closed before the connection is used again.

It automatically handles tracing.
*/
func (conn *connection) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: QueryRow", humanized))
	defer span.End()

	br := conn.client.SendBatch(ctx, batch)
	setDefaultAttributes(span, conn.config)
	setBatchAttributes(span, batch)

	return br
}

/*
WaitForNotification waits for a PostgreSQL notification.

It automatically handles tracing and error recording.
*/
func (conn *connection) WaitForNotification(ctx context.Context) (*pgconn.Notification, error) {
	ctx, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: WaitForNotification", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to wait for notification", err)
		}
	}()

	notif, err := conn.client.WaitForNotification(ctx)
	setDefaultAttributes(span, conn.config)

	return notif, err
}
