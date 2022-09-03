package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	cfg = cfg.Validate()

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
