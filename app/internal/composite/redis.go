package composite

import (
	"context"
	"fmt"

	"anilibrary-request-parser/app/internal/config"
	redisClient "anilibrary-request-parser/app/pkg/client/redis"
	"github.com/go-redis/redis/v9"
)

type RedisComposite struct {
	client *redis.Client
	cfg    config.Redis
}

func NewComposite(ctx context.Context, cfg config.Redis) (RedisComposite, error) {
	composite := RedisComposite{
		cfg: cfg,
	}

	client, err := redisClient.New(ctx, fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.Password)

	if err != nil {
		return composite, fmt.Errorf("while creating redis client : %w", err)
	}

	composite.client = client

	return composite, nil
}

func (c RedisComposite) Close() error {
	return c.client.Close()
}
