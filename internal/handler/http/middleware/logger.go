package middleware

import (
	"context"
	"errors"
	"net/http"

	"anilibrary-scraper/pkg/logger"
)

var ErrNoLogger = errors.New("logger not provided")

type ctxLogger struct{}

func Logger(log logger.Contract) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			withLogger := WithLogger(request.Context(), log)
			next.ServeHTTP(writer, request.WithContext(withLogger))
		})
	}
}

func WithLogger(ctx context.Context, contract logger.Contract) context.Context {
	return context.WithValue(ctx, ctxLogger{}, contract)
}

func GetLogger(ctx context.Context) logger.Contract {
	l, _ := ctx.Value(ctxLogger{}).(logger.Contract)
	return l
}

func MustGetLogger(ctx context.Context) logger.Contract {
	log, _ := ctx.Value(ctxLogger{}).(logger.Contract)
	if log == nil {
		panic(ErrNoLogger)
	}

	return log
}
