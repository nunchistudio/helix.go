package contextkey

/*
eventKeyIdentifier is the unique internal type to get/set an Event when interacting
with a Go context.
*/
type eventKeyIdentifier struct{}

/*
Event is the key identifier to get/set an Event when interacting with a Go
context.
*/
var Event eventKeyIdentifier
