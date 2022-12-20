package providers

import (
	"context"
	"fmt"
	"time"

	"anilibrary-scraper/pkg/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewJaegerTracerProvider(endpoint, serviceName, environment string, logger logging.Contract) (func(), error) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(endpoint),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("jaeger exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(environment),
		)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	closer := func() {
		logger.Info("closing jaeger")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := provider.Shutdown(ctx); err != nil {
			logger.Error("jaeger provider", logging.Error(err))
		}
	}

	return closer, nil
}
