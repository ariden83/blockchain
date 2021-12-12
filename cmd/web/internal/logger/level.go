package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// LevelsMap is a map to Logrus log levels
var LevelsMap = map[string]zapcore.Level{
	"EMERGENCY": zapcore.PanicLevel,
	"ALERT":     zapcore.FatalLevel,
	"CRITICAL":  zapcore.FatalLevel,
	"ERROR":     zapcore.ErrorLevel,
	"WARNING":   zapcore.WarnLevel,
	"NOTICE":    zapcore.InfoLevel,
	"INFO":      zapcore.InfoLevel,
	"DEBUG":     zapcore.DebugLevel,
}

// Level implements the UnmarshalText Interface to be able to be load into the config.
// It is mainly used for parsing configuration with github.com/kelseyhightower/envconfig.
type Level zapcore.Level

// UnmarshalText unmarshal a text to a level.
// Returns an error if no level match.
func (l *Level) UnmarshalText(text []byte) error {
	parsedLevel, found := LevelsMap[strings.ToUpper(string(text))]
	if !found {
		return fmt.Errorf("unable to parse level %s", text)
	}

	*l = Level(parsedLevel)
	return nil
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	return zapcore.Level(l).String()
}
