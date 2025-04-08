package event

import (
	"context"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := otel.Tracer("EventService").Start(ctx, "Send")
	defer span.End()

	span.AddEvent("Sending event")

	err := s.eventRepository.Send(ctx, model.Event{
		URL:       dto.URL,
		IP:        dto.IP,
		UserAgent: dto.UserAgent,
		Timestamp: dto.CreatedAt.Unix(),
	})
	if err != nil {
		span.SetStatus(codes.Error, "failed to send event")
		span.RecordError(err)

		return fmt.Errorf("sending event: %w", err)
	}

	span.SetStatus(codes.Ok, "event has been sent")

	return nil
}
