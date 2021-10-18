package logger

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// GELFEncoder is a gelf encoder.
// It will use an underlying json encoder to encode entries to the gelf format.
// See http://docs.graylog.org/en/2.4/pages/gelf.html.
type GELFEncoder struct {
	zapcore.Encoder
}

// Clone implements the encoder interface.
func (e GELFEncoder) Clone() zapcore.Encoder {
	return &GELFEncoder{e.Encoder.Clone()}
}

// EncodeEntry escape the keys following the gelf spec, then the underlying encoder will encode the entry.
func (e GELFEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	for i, f := range fields {
		fields[i].Key = escapeKey(f.Key)
	}
	return e.Encoder.EncodeEntry(entry, fields)
}

// NewGELFEncoder instanciate a new GELFEncoder.
func NewGELFEncoder() *GELFEncoder {
	return &GELFEncoder{
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			NameKey:        "_logger",
			LevelKey:       "level",
			CallerKey:      "_caller",
			MessageKey:     "short_message",
			StacktraceKey:  "full_message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeName:     zapcore.FullNameEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeLevel:    levelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		}),
	}
}

func escapeKey(key string) string {
	switch key {
	case "id":
		return "__id"
	case "version", "host", "short_message", "full_message", "timestamp", "level":
		return key
	}

	if len(key) == 0 {
		return "_"
	}

	if key[0] == '_' {
		return key
	}

	return "_" + key
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendInt(7)
	case zapcore.InfoLevel:
		enc.AppendInt(6)
	case zapcore.WarnLevel:
		enc.AppendInt(4)
	case zapcore.ErrorLevel:
		enc.AppendInt(3)
	case zapcore.DPanicLevel:
		enc.AppendInt(0)
	case zapcore.PanicLevel:
		enc.AppendInt(0)
	case zapcore.FatalLevel:
		enc.AppendInt(0)
	}
}
