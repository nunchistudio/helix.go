package contextkey

/*
spanKeyIdentifier is the unique internal type to get/set a custom OpenTelemetry
Span when interacting with a Go context.
*/
type spanKeyIdentifier struct{}

/*
Span is the key identifier to get/set a custom OpenTelemetry Span when interacting
with a Go context.
*/
var Span spanKeyIdentifier
