package middleware

import (
	"context"
	"errors"
	"net/http"

	"anilibrary-scraper/pkg/logging"
)

var ErrNoLogger = errors.New("logger not provided")

type ctxLogger struct{}

func Logger(log logging.Contract) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			withLogger := WithLogger(request.Context(), log)
			next.ServeHTTP(writer, request.WithContext(withLogger))
		})
	}
}

func WithLogger(ctx context.Context, contract logging.Contract) context.Context {
	return context.WithValue(ctx, ctxLogger{}, contract)
}

func GetLogger(ctx context.Context) logging.Contract {
	logger, _ := ctx.Value(ctxLogger{}).(logging.Contract)
	return logger
}

func MustGetLogger(ctx context.Context) logging.Contract {
	logger, _ := ctx.Value(ctxLogger{}).(logging.Contract)
	if logger == nil {
		panic(ErrNoLogger)
	}

	return logger
}
