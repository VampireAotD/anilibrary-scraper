package api

import (
	"context"
	"net"
	"strconv"
	"time"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 15 * time.Second
)

func NewServer(cfg config.HTTP, lifecycle fx.Lifecycle) fiber.Router {
	app := fiber.New(fiber.Config{
		AppName:      "Anilibrary API Server",
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logging.Get(),
	}))
	app.Use(recover.New())

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.MainPort))

			logging.Get().Info("HTTP server started at", zap.String("addr", addr))

			go func() {
				if err := app.Listen(addr); err != nil {
					logging.Get().Error("HTTP server failed to start", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})

	return app
}
