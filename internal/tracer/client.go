package tracer

import (
	"context"
	"os"
	"time"

	"go.nunchi.studio/helix/internal/cloudprovider"
	_ "go.nunchi.studio/helix/internal/setup"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

/*
t is the creator of Spans.
*/
var t trace.Tracer = otel.Tracer("go.nunchi.studio/helix/telemetry/trace")

/*
Tracer returns the global tracer used in the service.
*/
func Tracer() trace.Tracer {
	return t
}

/*
exporter holds the global tracer exporter used in the service.
*/
var exporter *otlptrace.Exporter

/*
Exporter returns the global tracer exporter used in the service.
*/
func Exporter() *otlptrace.Exporter {
	return exporter
}

/*
init initializes the global tracer client, with the appropriate attributes given
by the detected cloud provider. Panics if anything goes wrong (this should never
happen).
*/
func init() {
	var err error

	// Set appropriate tracer configuration.
	ctx := context.Background()
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 15 * time.Second,
			MaxInterval:     10 * time.Minute,
			MaxElapsedTime:  24 * time.Hour,
		}),
	}

	// Set the OpenTelemetry traces' endpoint only if specified.
	if os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")))
	}

	// Create the trace exporter.
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	if err != nil {
		panic(err)
	}

	// Get trace attributes returned by the detected cloud provider.
	resources, err := resource.New(ctx, resource.WithAttributes(cloudprovider.Detected.TracerAttributes()...))
	if err != nil {
		panic(err)
	}

	// Set the global trace reporter.
	otel.SetTracerProvider(sdk.NewTracerProvider(
		sdk.WithResource(resources),
		sdk.WithBatcher(exporter),
	))

	// Set the global trace propagator.
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
}
