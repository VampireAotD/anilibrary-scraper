package providers

import (
	"context"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/config"
	"github.com/VampireAotD/anilibrary-scraper/pkg/logging"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/fx"
)

func NewKafkaProvider(lifecycle fx.Lifecycle, cfg config.Kafka) (*kafka.Writer, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("creating scram mechanism: %w", err)
	}

	transport := &kafka.Transport{
		SASL: mechanism,
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Address),
		Topic:                  cfg.Topic,
		Balancer:               &kafka.LeastBytes{},
		Transport:              transport,
		AllowAutoTopicCreation: true,
		BatchTimeout:           cfg.BatchTimeout,
	}

	logging.Get().Info("Connected to Kafka")

	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			logging.Get().Info("Closing Kafka connection")
			return writer.Close()
		},
	})

	return writer, nil
}
