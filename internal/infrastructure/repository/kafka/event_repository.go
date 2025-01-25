package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/segmentio/kafka-go"
)

type EventRepository struct {
	connection *kafka.Conn
}

func NewEventRepository(connection *kafka.Conn) EventRepository {
	return EventRepository{
		connection: connection,
	}
}

func (r EventRepository) Send(_ context.Context, event model.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal 'event' model for Kafka: %w", err)
	}

	_, err = r.connection.WriteMessages(kafka.Message{
		Value: bytes,
	})
	if err != nil {
		return fmt.Errorf("send event to Kafka: %w", err)
	}

	return nil
}
