package app

import (
	"anilibrary-scraper/pkg/logging"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	logger          logging.Contract
	redisConnection *redis.Client
}

func SetupDependencies(logger logging.Contract, redisConnection *redis.Client) Dependencies {
	return Dependencies{
		logger:          logger,
		redisConnection: redisConnection,
	}
}
