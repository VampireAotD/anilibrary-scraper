package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/segmentio/kafka-go"
)

type EventRepository struct {
	writer *kafka.Writer
}

func NewEventRepository(writer *kafka.Writer) EventRepository {
	return EventRepository{
		writer: writer,
	}
}

func (r EventRepository) Send(ctx context.Context, event model.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal 'event' model for Kafka: %w", err)
	}

	err = r.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
	if err != nil {
		return fmt.Errorf("send event to Kafka: %w", err)
	}

	return nil
}
