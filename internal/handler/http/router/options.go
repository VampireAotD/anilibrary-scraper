package router

import (
	"anilibrary-scraper/pkg/logging"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Option func(cfg *config)

func WithLogger(logger logging.Contract) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

func WithRedisConnection(connection *redis.Client) Option {
	return func(cfg *config) {
		cfg.redisConnection = connection
	}
}

func WithKafkaConnection(connection *kafka.Conn) Option {
	return func(cfg *config) {
		cfg.kafkaConnection = connection
	}
}

func WithProfilingRoutes(enable bool) Option {
	return func(cfg *config) {
		cfg.enableProfiling = enable
	}
}
