package logging

import (
	"io"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Contract interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Named(s string) *Logger
	With(fields ...Field) *Logger
	Sync() error
}

type Logger struct {
	base *zap.Logger
}

func NewLogger(console io.Writer, files ...io.Writer) *Logger {
	pe := zap.NewProductionEncoderConfig()

	// file
	pe.CallerKey = "file"
	pe.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.CapitalString())
	}
	fileEncoder := zapcore.NewJSONEncoder(ecszap.ECSCompatibleEncoderConfig(pe))

	// console
	pe.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("|")
		encoder.AppendString(level.CapitalString())
		encoder.AppendString("|")
	}
	pe.ConsoleSeparator = " "
	pe.EncodeName = func(s string, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(s)
		encoder.AppendString("|")
	}
	consoleEncoder := zapcore.NewConsoleEncoder(ecszap.ECSCompatibleEncoderConfig(pe))

	cores := make([]zapcore.Core, len(files)+1)

	// console
	cores[0] = zapcore.NewCore(consoleEncoder,
		zapcore.AddSync(console),
		zap.DebugLevel,
	)

	for i := range files {
		cores[i+1] = zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(files[i]),
			zap.DebugLevel,
		)
	}

	return &Logger{base: zap.New(
		ecszap.WrapCore(zapcore.NewTee(cores...)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.base.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.base.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.base.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.base.Error(msg, fields...)
}

func (l *Logger) Named(s string) *Logger {
	if l.base == nil {
		return l
	}

	l.base = l.base.Named(s)
	return l
}

func (l *Logger) With(fields ...Field) *Logger {
	if l.base == nil {
		return l
	}

	l.base = l.base.With(fields...)
	return l
}

func (l *Logger) Sync() error {
	if l.base == nil {
		return nil
	}

	return l.base.Sync()
}
