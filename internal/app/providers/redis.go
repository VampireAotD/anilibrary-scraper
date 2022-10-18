package providers

import (
	"context"
	"fmt"
	"runtime"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-redis/redis/v9"
)

func NewRedisProvider(cfg config.Redis, log logger.Contract) (*redis.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.PoolTimeout)
	defer cancel()

	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10 * runtime.GOMAXPROCS(0)
	}

	if cfg.IdleSize <= 0 {
		cfg.IdleSize = cfg.PoolSize / 4
	}

	opts := &redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
		MinIdleConns: cfg.IdleSize,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("ping : %w", err)
	}

	cleanup := func() {
		log.Info("closing redis connection")
		if err := client.Close(); err != nil {
			log.Error("redis close", logger.Error(err))
		}
	}

	return client, cleanup, nil
}
