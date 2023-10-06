package logging

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	lvl          = zap.NewAtomicLevel()
	globalLogger *zap.Logger
)

// Get will return instance of configured logger using New. If none was configured - default logger will be provided.
func Get() *zap.Logger {
	if globalLogger == nil {
		globalLogger = New()
	}

	return globalLogger
}

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

	consoleEncoder := defaultConsoleEncoder(encCfg)
	fileEncoder := defaultFileEncoder(encCfg)

	cores := make([]zapcore.Core, 0, len(cfg.logFiles)+1)

	for i := range cfg.logFiles {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(zapcore.AddSync(cfg.logFiles[i])), lvl))
	}

	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(zapcore.AddSync(cfg.output)), lvl))

	tee := zapcore.NewTee(cores...)

	if cfg.ecsCompatible {
		tee = ecszap.WrapCore(tee)
	}

	globalLogger = zap.New(tee, zap.AddCaller())

	return globalLogger
}

func defaultConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	cfg.ConsoleSeparator = " "

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("|")
		encoder.AppendString(level.CapitalString())
		encoder.AppendString("|")
	}
	cfg.EncodeName = func(s string, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(s)
		encoder.AppendString("|")
	}

	return zapcore.NewConsoleEncoder(cfg)
}

func defaultFileEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	cfg.CallerKey = "file"
	cfg.TimeKey = "timestamp"
	cfg.MessageKey = "message"

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.CapitalString())
	}

	return zapcore.NewJSONEncoder(cfg)
}
