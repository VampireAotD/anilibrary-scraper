package logging

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New will create a new Logger instance with predefined encoders for log files and console output.
// Default output is os.Stdout and default level is InfoLevel.
func New(options ...Option) *Logger {
	cfg := &config{
		output: os.Stdout,
		level:  InfoLevel,
	}

	for i := range options {
		options[i](cfg)
	}

	encCfg := zap.NewProductionEncoderConfig()

	if cfg.ecsCompatible {
		encCfg = ecszap.ECSCompatibleEncoderConfig(encCfg)
	}

	core := zapcore.NewCore(
		resolveEncoder(encCfg, cfg.jsonEncoder),
		zapcore.Lock(zapcore.AddSync(cfg.output)),
		cfg.level,
	)

	if cfg.ecsCompatible {
		core = ecszap.WrapCore(core)
	}

	logger := zap.New(core, zap.AddCaller())

	if cfg.asDefault {
		SetDefault(logger)
	}

	return logger
}

func SetDefault(logger *Logger) {
	zap.ReplaceGlobals(logger)
}

// Get will return instance of configured logger using New.
// If logger wasn't configured - default logger will be provided.
func Get() *Logger {
	return zap.L()
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
