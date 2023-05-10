package providers

import (
	"context"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/segmentio/kafka-go"
)

func NewKafkaProvider(cfg config.Kafka, logger logging.Contract) (*kafka.Conn, func(), error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Address, cfg.Topic, cfg.Partition)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		logger.Info("closing kafka connection")

		if err := conn.Close(); err != nil {
			logger.Error("kafka provider", logging.Error(err))
		}
	}

	return conn, cleanup, nil
}
