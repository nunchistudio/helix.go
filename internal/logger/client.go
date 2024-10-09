package logger

import (
	"os"

	"go.nunchi.studio/helix/internal/orchestrator"
	_ "go.nunchi.studio/helix/internal/setup"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
client holds the global logger client used in the service.
*/
var client *zap.Logger

/*
Logger returns the global logger client used in the service.
*/
func Logger() *zap.Logger {
	return client
}

/*
init initializes the global logger client, with the appropriate fields given by
the detected orchestrator. Panics if anything goes wrong (this should never
happen).
*/
func init() {
	var err error

	// Set the appropriate log level given the "ENVIRONMENT" environment variable.
	var level zapcore.Level
	switch os.Getenv("ENVIRONMENT") {
	case "local", "localhost", "dev", "development":
		level = zapcore.DebugLevel
	default:
		level = zapcore.InfoLevel
	}

	// Set appropriate logger configuration.
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Override some keys and time encoding for better consistency.
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	// Get log fields returned by the detected orchestrator.
	fields := orchestrator.Detected.LoggerFields()

	// Finally, try to create the global logger client.
	client, err = cfg.Build(zap.Fields(fields...))
	if err != nil {
		panic(err)
	}
}
