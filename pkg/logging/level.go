package logging

import "go.uber.org/zap/zapcore"

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

func SetLevel(level string) {
	switch level {
	case InfoLevel:
		lvl.SetLevel(zapcore.InfoLevel)
	case ErrorLevel:
		lvl.SetLevel(zapcore.ErrorLevel)
	case PanicLevel:
		lvl.SetLevel(zapcore.PanicLevel)
	case WarnLevel:
		lvl.SetLevel(zapcore.WarnLevel)
	case FatalLevel:
		lvl.SetLevel(zapcore.FatalLevel)
	default:
		lvl.SetLevel(zapcore.DebugLevel)
	}
}
