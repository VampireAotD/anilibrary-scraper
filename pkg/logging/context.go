package logging

import (
	"context"
)

type loggerCxtKey struct{}

func ContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCxtKey{}, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerCxtKey{}).(*Logger); ok {
		return logger
	}

	return Get()
}
