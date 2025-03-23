package providers

import (
	"context"
	"fmt"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/config"
	"github.com/VampireAotD/anilibrary-scraper/pkg/logging"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"go.uber.org/fx"
)

func NewKafkaProvider(lifecycle fx.Lifecycle, cfg config.Kafka) (*kgo.Client, error) {
	auth := scram.Auth{
		User: cfg.Username,
		Pass: cfg.Password,
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Address),
		kgo.SASL(auth.AsSha512Mechanism()),
		kgo.DefaultProduceTopic(cfg.Topic),
		kgo.ProducerBatchMaxBytes(1024*1024),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating Kafka client: %w", err)
	}

	logging.Get().Info("Connected to Kafka")

	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			logging.Get().Info("Closing Kafka connection")
			client.Close()
			return nil
		},
	})

	return client, nil
}
