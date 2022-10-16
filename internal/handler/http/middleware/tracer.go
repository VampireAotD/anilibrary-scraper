package middleware

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var ErrNoTracer = errors.New("tracer not provided")

type ctxTracer struct{}

func Tracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		withTracer := WithTracer(request.Context())
		next.ServeHTTP(writer, request.WithContext(withTracer))
	})
}

func WithTracer(ctx context.Context) context.Context {
	tracer := otel.Tracer("http")
	return context.WithValue(ctx, ctxTracer{}, tracer)
}

func MustGetTracer(ctx context.Context) trace.Tracer {
	tracer, _ := ctx.Value(ctxTracer{}).(trace.Tracer)
	if tracer == nil {
		panic(ErrNoTracer)
	}

	return tracer
}
