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
	eventRepository Repository
}

func NewService(eventRepository Repository) Service {
	return Service{
		eventRepository: eventRepository,
	}
}

func (s Service) Send(ctx context.Context, dto DTO) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("EventService").Start(ctx, "Send")
	defer span.End()

	span.AddEvent("Sending event to ClickHouse")

	return s.eventRepository.Send(ctx, model.Event{
		URL:       dto.URL,
		IP:        dto.IP,
		UserAgent: dto.UserAgent,
		Timestamp: time.Now().Unix(),
	})
}
