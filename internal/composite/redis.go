package composite

import (
	"context"
	"fmt"

	"anilibrary-request-parser/internal/config"
	redisClient "anilibrary-request-parser/pkg/client/redis"
	"github.com/go-redis/redis/v9"
)

type RedisComposite struct {
	Client *redis.Client
}

func NewRedisComposite(cfg config.Redis) (RedisComposite, error) {
	var composite RedisComposite

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PoolTimeout)

	defer cancel()

	client, err := redisClient.New(ctx, redisClient.Config{
		Address:     cfg.Address,
		Password:    cfg.Password,
		PoolTimeout: cfg.PoolTimeout,
		PoolSize:    cfg.PoolSize,
		IdleSize:    cfg.IdleSize,
	})

	if err != nil {
		return composite, fmt.Errorf("while creating redis client : %w", err)
	}

	composite.Client = client

	return composite, nil
}

func (c RedisComposite) Close() error {
	return c.Client.Close()
}
