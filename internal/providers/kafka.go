package providers

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

func NewKafkaProvider(lifecycle fx.Lifecycle, cfg config.Kafka) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Address, cfg.Topic, cfg.Partition)
	if err != nil {
		return nil, fmt.Errorf("connecting to kafka: %w", err)
	}

	logging.Get().Info("Connected to Kafka")

	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			logging.Get().Info("Closing Kafka connection")

			return conn.Close()
		},
	})

	return conn, nil
}
