package event

import (
	"context"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/service"
	"go.opentelemetry.io/otel/trace"
)

var _ service.EventService = (*Service)(nil)

type Service struct {
	kafkaRepository repository.EventRepository
}

func NewService(kafkaRepository repository.EventRepository) Service {
	return Service{
		kafkaRepository: kafkaRepository,
	}
}

func (s Service) Send(ctx context.Context, url string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("EventService").Start(ctx, "Send")
	defer span.End()

	span.AddEvent("Sending event to Clickhouse")

	return s.kafkaRepository.Send(ctx, &entity.Event{
		URL:  url,
		Date: time.Now().Unix(),
	})
}
