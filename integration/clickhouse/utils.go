package clickhouse

import (
	"fmt"
	"unicode"

	"go.nunchi.studio/helix/telemetry/trace"
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
normalizeErrorMessage normalizes an error returned by the ClickHouse client to
match the format of helix.go. This is only used inside Start and Close for a
better readability in the terminal. Otherwise, functions return native ClickHouse
errors.

Example:

	"dial tcp 127.0.0.1:8123: connect: connection refused"

Becomes:

	"Dial tcp 127.0.0.1:8123: connect: connection refused"
*/
func normalizeErrorMessage(err error) string {
	var msg string = err.Error()
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
