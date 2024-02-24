package logging

import (
	"os"
	"sync"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	lvl          = zap.NewAtomicLevel()
	globalLogger *zap.Logger
	once         sync.Once
)

// New will create a new instance of *zap.Logger with predefined encoders for log files and console output.
// Default output is os.Stdout.
func New(options ...Option) *zap.Logger {
	cfg := &config{
		output: os.Stdout,
	}

	for i := range options {
		options[i](cfg)
	}

	encCfg := zap.NewProductionEncoderConfig()

	if cfg.ecsCompatible {
		encCfg = ecszap.ECSCompatibleEncoderConfig(encCfg)
	}

	core := zapcore.NewCore(resolveEncoder(encCfg, cfg.jsonEncoder), zapcore.Lock(zapcore.AddSync(cfg.output)), lvl)

	if cfg.ecsCompatible {
		core = ecszap.WrapCore(core)
	}

	globalLogger = zap.New(core, zap.AddCaller())

	return globalLogger
}

// Get will return instance of configured logger using New. If none was configured - default logger will be provided.
func Get() *zap.Logger {
	once.Do(func() {
		if globalLogger == nil {
			globalLogger = New()
		}
	})

	return globalLogger
}

func resolveEncoder(cfg zapcore.EncoderConfig, encodeJSON bool) zapcore.Encoder {
	cfg.CallerKey = "file"
	cfg.TimeKey = "timestamp"
	cfg.MessageKey = "message"

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.CapitalString())
	}

	if encodeJSON {
		return zapcore.NewJSONEncoder(cfg)
	}

	return zapcore.NewConsoleEncoder(cfg)
}
