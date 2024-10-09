package openfeature

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.nunchi.studio/helix/errorstack"
	"go.nunchi.studio/helix/event"
	"go.nunchi.studio/helix/internal/logger"
	"go.nunchi.studio/helix/internal/orchestrator"
	"go.nunchi.studio/helix/service"

	"github.com/go-logr/zapr"
	gofeatureflaginprocess "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag-in-process/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

/*
OpenFeature exposes an opinionated way to interact with the OpenFeature specification
and GO Feature Flag provider, by bringing automatic distributed tracing as well
as error recording within traces.
*/
type OpenFeature interface {
	EvaluateString(ctx context.Context, flag string, defaultValue string, target string) (openfeature.StringEvaluationDetails, error)
	EvaluateBoolean(ctx context.Context, flag string, defaultValue bool, target string) (openfeature.BooleanEvaluationDetails, error)
	EvaluateInteger(ctx context.Context, flag string, defaultValue int64, target string) (openfeature.IntEvaluationDetails, error)
	EvaluateFloat(ctx context.Context, flag string, defaultValue float64, target string) (openfeature.FloatEvaluationDetails, error)
}

/*
connection represents the openfeature integration. It respects the
integration.Integration and OpenFeature interfaces.
*/
type connection struct {

	// config holds the Config initially passed when creating a new OpenFeature
	// client.
	config *Config

	// client is the connection made with the OpenFeature client.
	client *openfeature.Client
}

/*
Init tries to create an OpenFeature client with the GO Feature Flag provider
given the Config. Returns an error if Config is not valid or if the initialization
failed.
*/
func Init(cfg Config) (OpenFeature, error) {

	// No need to continue if Config is not valid.
	err := cfg.sanitize()
	if err != nil {
		return nil, err
	}

	// Start to build an error stack, so we can add validations as we go.
	stack := errorstack.New("Failed to initialize integration", errorstack.WithIntegration(identifier))
	conn := &connection{
		config: &cfg,
	}

	// Create a file retriever for each path set in Config.
	var retrievers []retriever.Retriever
	for _, path := range cfg.Paths {
		retrievers = append(retrievers, &fileretriever.Retriever{
			Path: path,
		})
	}

	// Set the OpenFeature options for the GO Feature Flag provider.
	var opts = gofeatureflaginprocess.ProviderOptions{
		GOFeatureFlagConfig: &ffclient.Config{
			Environment: os.Getenv("ENVIRONMENT"),
			Retrievers:  retrievers,
		},
	}

	// Try to create the OpenFeature provider.
	provider, err := gofeatureflaginprocess.NewProvider(opts)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return nil, stack
	}

	// Set the default OpenFeature provider.
	err = openfeature.SetProviderAndWait(provider)
	if err != nil {
		stack.WithValidations(errorstack.Validation{
			Message: err.Error(),
		})

		return nil, stack
	}

	// Finally, create the OpenFeature client with the appropriate service name
	// and global logger.
	conn.client = openfeature.NewClient(orchestrator.Detected.Service())
	conn.client.WithLogger(zapr.NewLogger(logger.Logger()))

	// Try to attach the integration to the service.
	if err := service.Attach(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

/*
EvaluateString evaluates a string value for a given flag against a specific target.
The evaluation context is the JSON marshaled event.Event found in context.

It automatically handles tracing (new event in current span) and error recording.
*/
func (conn *connection) EvaluateString(ctx context.Context, flag string, defaultValue string, target string) (openfeature.StringEvaluationDetails, error) {
	var details openfeature.StringEvaluationDetails
	var err error

	span := trace.SpanFromContext(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to evaluate string target")
		}
	}()

	var mapped map[string]any
	e, found := event.EventFromContext(ctx)
	if found {
		b, err := json.Marshal(e)
		if err != nil {
			return details, err
		}

		err = json.Unmarshal(b, &mapped)
		if err != nil {
			return details, err
		}
	}

	details, err = conn.client.StringValueDetails(ctx, flag, defaultValue, openfeature.NewEvaluationContext(target, mapped))
	if err != nil {
		return details, err
	}

	span.AddEvent(fmt.Sprintf("%s.evaluate", identifier),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.flag", identifier), flag)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.default_value", identifier), defaultValue)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.target", identifier), target)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.variant", identifier), details.Variant)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.value", identifier), details.Value)),
	)

	return details, err
}

/*
EvaluateBoolean evaluates a boolean value for a given flag against a specific
target. The evaluation context is the JSON marshaled event.Event found in context.

It automatically handles tracing (new event in current span) and error recording.
*/
func (conn *connection) EvaluateBoolean(ctx context.Context, flag string, defaultValue bool, target string) (openfeature.BooleanEvaluationDetails, error) {
	var details openfeature.BooleanEvaluationDetails
	var err error

	span := trace.SpanFromContext(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to evaluate boolean target")
		}
	}()

	var mapped map[string]any
	e, found := event.EventFromContext(ctx)
	if found {
		b, err := json.Marshal(e)
		if err != nil {
			return details, err
		}

		err = json.Unmarshal(b, &mapped)
		if err != nil {
			return details, err
		}
	}

	details, err = conn.client.BooleanValueDetails(ctx, flag, defaultValue, openfeature.NewEvaluationContext(target, mapped))
	if err != nil {
		return details, err
	}

	span.AddEvent(fmt.Sprintf("%s.evaluate", identifier),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.flag", identifier), flag)),
		trace.WithAttributes(attribute.Bool(fmt.Sprintf("%s.default_value", identifier), defaultValue)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.target", identifier), target)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.variant", identifier), details.Variant)),
		trace.WithAttributes(attribute.Bool(fmt.Sprintf("%s.value", identifier), details.Value)),
	)

	return details, err
}

/*
EvaluateInteger evaluates an integer value for a given flag against a specific
target. The evaluation context is the JSON marshaled event.Event found in context.

It automatically handles tracing (new event in current span) and error recording.
*/
func (conn *connection) EvaluateInteger(ctx context.Context, flag string, defaultValue int64, target string) (openfeature.IntEvaluationDetails, error) {
	var details openfeature.IntEvaluationDetails
	var err error

	span := trace.SpanFromContext(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to evaluate integer target")
		}
	}()

	var mapped map[string]any
	e, found := event.EventFromContext(ctx)
	if found {
		b, err := json.Marshal(e)
		if err != nil {
			return details, err
		}

		err = json.Unmarshal(b, &mapped)
		if err != nil {
			return details, err
		}
	}

	details, err = conn.client.IntValueDetails(ctx, flag, defaultValue, openfeature.NewEvaluationContext(target, mapped))
	if err != nil {
		return details, err
	}

	span.AddEvent(fmt.Sprintf("%s.evaluate", identifier),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.flag", identifier), flag)),
		trace.WithAttributes(attribute.Int64(fmt.Sprintf("%s.default_value", identifier), defaultValue)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.target", identifier), target)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.variant", identifier), details.Variant)),
		trace.WithAttributes(attribute.Int64(fmt.Sprintf("%s.value", identifier), details.Value)),
	)

	return details, err
}

/*
EvaluateFloat evaluates a float value for a given flag against a specific target.
The evaluation context is the JSON marshaled event.Event found in context.

It automatically handles tracing (new event in current span) and error recording.
*/
func (conn *connection) EvaluateFloat(ctx context.Context, flag string, defaultValue float64, target string) (openfeature.FloatEvaluationDetails, error) {
	var details openfeature.FloatEvaluationDetails
	var err error

	span := trace.SpanFromContext(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to evaluate float target")
		}
	}()

	var mapped map[string]any
	e, found := event.EventFromContext(ctx)
	if found {
		b, err := json.Marshal(e)
		if err != nil {
			return details, err
		}

		err = json.Unmarshal(b, &mapped)
		if err != nil {
			return details, err
		}
	}

	details, err = conn.client.FloatValueDetails(ctx, flag, defaultValue, openfeature.NewEvaluationContext(target, mapped))
	if err != nil {
		return details, err
	}

	span.AddEvent(fmt.Sprintf("%s.evaluate", identifier),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.flag", identifier), flag)),
		trace.WithAttributes(attribute.Float64(fmt.Sprintf("%s.default_value", identifier), defaultValue)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.target", identifier), target)),
		trace.WithAttributes(attribute.String(fmt.Sprintf("%s.variant", identifier), details.Variant)),
		trace.WithAttributes(attribute.Float64(fmt.Sprintf("%s.value", identifier), details.Value)),
	)

	return details, err
}
