package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func New(ctx context.Context, addr, password string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping : %w", err)
	}

	return client, nil
}
