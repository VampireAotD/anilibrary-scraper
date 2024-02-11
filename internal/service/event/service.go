package event

import (
	"context"
	"time"

	"anilibrary-scraper/internal/repository/model"

	"go.opentelemetry.io/otel/trace"
)

//go:generate mockgen -source=service.go -destination=./mocks.go -package=event
type Repository interface {
	// Send method sends event data
	Send(ctx context.Context, event model.Event) error
}

type Service struct {
	kafkaRepository Repository
}

func NewService(kafkaRepository Repository) Service {
	return Service{
		kafkaRepository: kafkaRepository,
	}
}

func (s Service) Send(ctx context.Context, url string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("EventService").Start(ctx, "Send")
	defer span.End()

	span.AddEvent("Sending event to Clickhouse")

	return s.kafkaRepository.Send(ctx, model.Event{
		URL:  url,
		Date: time.Now().Unix(),
	})
}
