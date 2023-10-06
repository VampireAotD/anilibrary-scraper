package middleware

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var ErrNoLogger = errors.New("logger not provided")

type ctxLogger struct{}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			withLogger := WithLogger(request.Context(), logger)
			next.ServeHTTP(writer, request.WithContext(withLogger))
		})
	}
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(ctxLogger{}).(*zap.Logger)
	return logger
}

func MustGetLogger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(ctxLogger{}).(*zap.Logger)
	if logger == nil {
		panic(ErrNoLogger)
	}

	return logger
}
