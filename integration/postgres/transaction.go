package postgres

import (
	"context"
	"fmt"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*
transaction implements the Tx interface and allows to wrap the PostgreSQL
transaction functions for automatic tracing and error recording.
*/
type transaction struct {
	config *Config
	client pgx.Tx
}

/*
Tx represents a database transaction.
*/
type Tx interface {
	Begin(ctx context.Context) (Tx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, query string, args ...any) (commandTag pgconn.CommandTag, err error)
	Prepare(ctx context.Context, id string, query string) (*pgconn.StatementDescription, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	LargeObjects(ctx context.Context) pgx.LargeObjects
}

/*
Begin starts a pseudo nested transaction implemented with a savepoint.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Begin(ctx context.Context) (Tx, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Begin", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to begin transaction", err)
		}
	}()

	subtx, err := tx.client.Begin(ctx)
	setDefaultAttributes(span, tx.config)

	sub := &transaction{
		config: tx.config,
		client: subtx,
	}

	return sub, err
}

/*
Commit commits the transaction.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Commit(ctx context.Context) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Commit", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to commit transaction", err)
		}
	}()

	err = tx.client.Commit(ctx)
	setDefaultAttributes(span, tx.config)

	return err
}

/*
Rollback rolls back the transaction. Rollback will return pgx.ErrTxClosed if the
Tx is already closed, but is otherwise safe to call multiple times. Hence, a
defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error
condition.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Rollback(ctx context.Context) error {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Rollback", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to rollback transaction", err)
		}
	}()

	err = tx.client.Rollback(ctx)
	setDefaultAttributes(span, tx.config)

	return err
}

/*
Exec delegates to the underlying connection.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Exec", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to execute query", err)
		}
	}()

	stmt, err := tx.client.Exec(ctx, query, args...)
	setDefaultAttributes(span, tx.config)
	setTransactionQueryAttributes(span, query)

	return stmt, err
}

/*
Prepare delegates to the underlying connection.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Prepare(ctx context.Context, id string, query string) (*pgconn.StatementDescription, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / Prepare", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to prepare statement", err)
		}
	}()

	stmt, err := tx.client.Prepare(ctx, id, query)
	setDefaultAttributes(span, tx.config)
	setTransactionQueryAttributes(span, query)

	return stmt, err
}

/*
Query delegates to the underlying connection.

It automatically handles tracing and error recording.
*/
func (tx *transaction) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / QueryRows", humanized))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.RecordError("failed to query rows", err)
		}
	}()

	rows, err := tx.client.Query(ctx, query, args...)
	setDefaultAttributes(span, tx.config)
	setTransactionQueryAttributes(span, query)

	return rows, err
}

/*
QueryRow delegates to the underlying connection.

It automatically handles tracing.
*/
func (tx *transaction) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / QueryRow", humanized))
	defer span.End()

	row := tx.client.QueryRow(ctx, query, args...)
	setDefaultAttributes(span, tx.config)
	setTransactionQueryAttributes(span, query)

	return row
}

/*
SendBatch delegates to the underlying connection.

It automatically handles tracing.
*/
func (tx *transaction) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	_, span := trace.Start(ctx, trace.SpanKindClient, fmt.Sprintf("%s: Transaction / SendBatch", humanized))
	defer span.End()

	res := tx.client.SendBatch(ctx, batch)
	setDefaultAttributes(span, tx.config)
	setTransactionBatchAttributes(span, batch)

	return res
}

/*
LargeObjects returns a pgx.LargeObjects instance for the transaction.
*/
func (tx *transaction) LargeObjects(ctx context.Context) pgx.LargeObjects {
	return tx.client.LargeObjects()
}
