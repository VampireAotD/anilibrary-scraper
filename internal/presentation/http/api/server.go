package api

import (
	"context"
	"net"
	"strconv"
	"time"

	"anilibrary-scraper/internal/infrastructure/config"
	"anilibrary-scraper/internal/presentation/http/middleware"
	"anilibrary-scraper/pkg/logging"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

const (
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 15 * time.Second
)

func NewServer(lifecycle fx.Lifecycle, cfg config.HTTP) fiber.Router {
	app := fiber.New(fiber.Config{
		AppName:      "Anilibrary API Server",
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logging.Get(),
	}))
	app.Use(middleware.NewLogger(logging.Get()))
	app.Use(recover.New())

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.MainPort))

			logging.Get().Info("HTTP server started at", logging.String("addr", addr))

			go func() {
				if err := app.Listen(addr); err != nil {
					logging.Get().Error("HTTP server failed to start", logging.Error(err))
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
