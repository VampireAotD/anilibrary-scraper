package composite

import (
	"context"
	"fmt"

	"anilibrary-request-parser/app/internal/config"
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

	client, err := composite.newClient(ctx)

	if err != nil {
		return composite, fmt.Errorf("while creating redis client : %w", err)
	}

	composite.client = client

	return composite, nil
}

func (c RedisComposite) Close() error {
	return c.client.Close()
}

func (c *RedisComposite) newClient(ctx context.Context) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		Password: c.cfg.Password,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping : %w", err)
	}

	return client, nil
}
