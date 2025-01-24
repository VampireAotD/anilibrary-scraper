package providers

import (
	"context"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/config"
	"github.com/VampireAotD/anilibrary-scraper/pkg/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/fx"
)

func NewTraceProvider(lifecycle fx.Lifecycle, cfg config.App) error {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return fmt.Errorf("trace exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Name),
			semconv.DeploymentEnvironmentKey.String(string(cfg.Env)),
		)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logging.Get().Info("Configured trace provider")

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logging.Get().Info("Closing tracer")

			return provider.Shutdown(ctx)
		},
	})

	return nil
}
