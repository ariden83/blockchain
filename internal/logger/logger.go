package logger

import (
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Init(cfg config.Log) *zap.Logger {

	opt := dir.Options{
		Dir:      cfg.Path,
		FileMode: os.FileMode(0705),
	}

	if !dir.DirExists(opt) {
		if err := dir.CreateDir(opt); err != nil {
			panic(err)
		}
	}

	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	wsError, err := os.OpenFile(opt.Dir+"/errors.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	wsAllLogs, err := os.OpenFile(opt.Dir+"/logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	// consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileProdEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	consoleGELFEncoder := NewGELFEncoder()

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleGELFEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleGELFEncoder, consoleDebugging, lowPriority),

		zapcore.NewCore(fileProdEncoder, wsError, highPriority),
		zapcore.NewCore(fileEncoder, wsAllLogs, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	//logger := zap.New(core)
	c := zap.NewProductionConfig()
	CLILevel := Level(LevelsMap[cfg.CLILevel])

	c.Level = zap.NewAtomicLevelAt(zapcore.Level(CLILevel))
	logger, err := c.Build()
	if err != nil {
		panic(err)
	}

	logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core
	}))

	return logger
}

func InitLight(cfg config.Log) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	// consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleGELFEncoder := NewGELFEncoder()

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleGELFEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleGELFEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)

	return logger
}
