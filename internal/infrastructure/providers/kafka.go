package providers

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/infrastructure/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/fx"
)

func NewKafkaProvider(lifecycle fx.Lifecycle, cfg config.Kafka) (*kafka.Conn, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("creating scram mechanism: %w", err)
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	conn, err := dialer.DialLeader(context.Background(), "tcp", cfg.Address, cfg.Topic, cfg.Partition)
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
