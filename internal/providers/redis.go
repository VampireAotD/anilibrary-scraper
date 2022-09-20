package providers

import (
	"context"
	"fmt"
	"runtime"

	"anilibrary-request-parser/internal/config"
	"github.com/go-redis/redis/v9"
)

func NewRedisProvider(cfg config.Redis) (*redis.Client, error) {
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
		return nil, fmt.Errorf("ping : %w", err)
	}

	return client, nil
}
