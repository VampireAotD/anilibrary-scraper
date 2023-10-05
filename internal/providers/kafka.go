package providers

import (
	"context"

	"anilibrary-scraper/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

func NewKafkaProvider(lifecycle fx.Lifecycle, cfg config.Kafka, logger logging.Contract) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Address, cfg.Topic, cfg.Partition)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing kafka connection")

			return conn.Close()
		},
	})

	return conn, nil
}
