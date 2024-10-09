package temporal

import (
	"context"

	"go.nunchi.studio/helix/telemetry/log"
)

/*
customloggeri mplements the Temporal's log.Logger interface, allowing to leverage
the global logger set by helix.go.
*/
type customlogger struct{}

/*
Debug logs a message at the debug level, using the global logger.
*/
func (l *customlogger) Debug(msg string, keyvals ...any) {
	log.Debug(context.Background(), msg)
}

/*
Info logs a message at the info level, using the global logger.
*/
func (l *customlogger) Info(msg string, keyvals ...any) {
	log.Info(context.Background(), msg)
}

/*
Warn logs a message at the warn level, using the global logger.
*/
func (l *customlogger) Warn(msg string, keyvals ...any) {
	log.Warn(context.Background(), msg)
}

/*
Error logs a message at the error level, using the global logger.
*/
func (l *customlogger) Error(msg string, keyvals ...any) {
	log.Error(context.Background(), msg)
}
