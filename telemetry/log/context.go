package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

/*
fromContextToZapFields tries to extract a trace from the context to add "trace_id"
and "span_id" fields to the log.
*/
func fromContextToZapFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field

	link := trace.LinkFromContext(ctx)
	if link.SpanContext.HasTraceID() {
		fields = append(fields, zapcore.Field{
			Key:    "trace_id",
			Type:   zapcore.StringType,
			String: link.SpanContext.TraceID().String(),
		})

		if link.SpanContext.HasSpanID() {
			fields = append(fields, zapcore.Field{
				Key:    "span_id",
				Type:   zapcore.StringType,
				String: link.SpanContext.SpanID().String(),
			})
		}
	}

	return fields
}
