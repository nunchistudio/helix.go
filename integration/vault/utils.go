package vault

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
		span.SetStringAttribute(fmt.Sprintf("%s.server.address", identifier), cfg.Address)
		span.SetStringAttribute(fmt.Sprintf("%s.agent.address", identifier), cfg.AgentAddress)
		span.SetStringAttribute(fmt.Sprintf("%s.namespace", identifier), cfg.Namespace)
	}
}

/*
setKeyValueAttributes sets Key-Value attributes to a trace span.
*/
func setKeyValueAttributes(span *trace.Span, mountpath string, secretpath string) {
	span.SetStringAttribute(fmt.Sprintf("%s.kv.mountpath", identifier), mountpath)
	span.SetStringAttribute(fmt.Sprintf("%s.kv.secretpath", identifier), secretpath)
}

/*
normalizeErrorMessage normalizes an error returned by the Vault client to match
the format of helix.go. This is only used inside Start and Close for a better
readability in the terminal. Otherwise, functions return native Vault errors.

Example:

	"dial tcp 127.0.0.1:8200: connect: connection refused"

Becomes:

	"Dial tcp 127.0.0.1:8200: connect: connection refused"
*/
func normalizeErrorMessage(err error) string {
	var msg string = err.Error()
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
