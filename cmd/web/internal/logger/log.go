package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger instanciate a new zap.Logger that will output to both console and graylog.
// If GraylogEndpoint == "", no data will be send to graylog.
func NewLogger(GraylogEndpoint string, GraylogLevel, CLILevel Level) (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.Level = zap.NewAtomicLevelAt(zapcore.Level(CLILevel))
	log, err := c.Build()
	if err != nil {
		return nil, err
	}

	if GraylogEndpoint == "" {
		return log, nil
	}

	graylogWriter, err := NewGELFWriter(GraylogEndpoint)
	if err != nil {
		return nil, err
	}

	log = log.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			c,
			zapcore.NewSampler(
				zapcore.NewCore(NewGELFEncoder(), graylogWriter, zap.NewAtomicLevelAt(zapcore.Level(GraylogLevel))),
				time.Second,
				100,
				100,
			),
		)
	}))

	// Add hostname if we can get it.
	if host, err := os.Hostname(); err == nil {
		log = log.With(zap.String("host", host))
	}

	return log, nil
}
