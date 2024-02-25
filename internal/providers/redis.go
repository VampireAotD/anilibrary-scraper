package providers

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func NewRedisProvider(lifecycle fx.Lifecycle, cfg config.Redis) (redis.UniversalClient, error) {
	opts := &redis.Options{
		Addr:        cfg.Address,
		Password:    cfg.Password,
		PoolTimeout: cfg.PoolTimeout,
	}

	client := redis.NewClient(opts)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("redis ping : %w", err)
			}

			if err := redisotel.InstrumentTracing(client); err != nil {
				return fmt.Errorf("redis tracing: %w", err)
			}

			logging.Get().Info("Connected to Redis")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logging.Get().Info("Closing Redis connection")

			return client.Close()
		},
	})

	return client, nil
}
