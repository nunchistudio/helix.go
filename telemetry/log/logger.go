package log

import (
	"context"

	"go.nunchi.studio/helix/internal/logger"
)

/*
Debug logs a message at the debug level. It tries to extract a trace from the
context to add "trace_id" and "span_id" fields to the log.
*/
func Debug(ctx context.Context, msg string) {
	logger.Logger().Debug(msg, fromContextToZapFields(ctx)...)
}

/*
Info logs a message at the info level. It tries to extract a trace from the
context to add "trace_id" and "span_id" fields to the log.
*/
func Info(ctx context.Context, msg string) {
	logger.Logger().Info(msg, fromContextToZapFields(ctx)...)
}

/*
Warn logs a message at the warn level. It tries to extract a trace from the
context to add "trace_id" and "span_id" fields to the log.
*/
func Warn(ctx context.Context, msg string) {
	logger.Logger().Warn(msg, fromContextToZapFields(ctx)...)
}

/*
Error logs a message at the error level. It tries to extract a trace from the
context to add "trace_id" and "span_id" fields to the log.
*/
func Error(ctx context.Context, msg string) {
	logger.Logger().Error(msg, fromContextToZapFields(ctx)...)
}

/*
Fatal logs a message at the fatal level. It tries to extract a trace from the
context to add "trace_id" and "span_id" fields to the log.

The logger then calls os.Exit(1).
*/
func Fatal(ctx context.Context, msg string) {
	logger.Logger().Fatal(msg, fromContextToZapFields(ctx)...)
}
