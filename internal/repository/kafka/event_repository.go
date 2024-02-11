package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"anilibrary-scraper/internal/repository/model"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type EventRepository struct {
	connection *kafka.Conn
}

func NewEventRepository(connection *kafka.Conn) EventRepository {
	return EventRepository{
		connection: connection,
	}
}

func (r EventRepository) Send(ctx context.Context, event model.Event) error {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("EventRepository").Start(ctx, "Send")
	defer span.End()

	bytes, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return err
	}

	_, err = r.connection.WriteMessages(kafka.Message{
		Value: bytes,
	})
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("sending event: %w", err)
	}

	return nil
}
