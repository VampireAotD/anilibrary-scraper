package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EventRepository struct {
	client *kgo.Client
}

func NewEventRepository(client *kgo.Client) EventRepository {
	return EventRepository{
		client: client,
	}
}

func (r EventRepository) Send(ctx context.Context, event model.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal 'event' model for Kafka: %w", err)
	}

	record := &kgo.Record{
		Value: bytes,
	}

	err = r.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("send event to Kafka: %w", err)
	}

	return nil
}
