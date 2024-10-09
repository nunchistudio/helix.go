package postgres

import (
	"fmt"
	"unicode"

	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/jackc/pgx/v5"
)

/*
setDefaultAttributes sets integration attributes to a trace span.
*/
func setDefaultAttributes(span *trace.Span, cfg *Config) {
	if cfg != nil {
		span.SetStringAttribute(fmt.Sprintf("%s.database", identifier), cfg.Database)
	}
}

/*
setQueryAttributes sets SQL query attributes to a trace span.
*/
func setQueryAttributes(span *trace.Span, query string) {
	span.SetStringAttribute(fmt.Sprintf("%s.query", identifier), query)
}

/*
setTransactionQueryAttributes sets SQL query attributes of a transaction to trace
span.
*/
func setTransactionQueryAttributes(span *trace.Span, query string) {
	span.SetStringAttribute(fmt.Sprintf("%s.transaction.query", identifier), query)
}

/*
setBatchAttributes sets batch attributes to a trace span.
*/
func setBatchAttributes(span *trace.Span, batch *pgx.Batch) {
	span.SetIntAttribute(fmt.Sprintf("%s.batch.length", identifier), int64(batch.Len()))
}

/*
setTransactionBatchAttributes sets batch attributes of a transaction to a trace
span.
*/
func setTransactionBatchAttributes(span *trace.Span, batch *pgx.Batch) {
	span.SetIntAttribute(fmt.Sprintf("%s.transaction.batch.length", identifier), int64(batch.Len()))
}

/*
normalizeErrorMessage normalizes an error returned by the PostgreSQL client to
match the format of helix.go. This is only used inside Start and Close for a
better readability in the terminal. Otherwise, functions return native PostgreSQL
errors.

Example:

	"failed to connect to `host=localhost user=postgres database=postgres`: ..."

Becomes:

	"Failed to connect to `host=localhost user=postgres database=postgres`: ..."
*/
func normalizeErrorMessage(err error) string {
	var msg string = err.Error()
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
