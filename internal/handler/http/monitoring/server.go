package monitoring

import (
	"context"
	"net"
	"strconv"

	"anilibrary-scraper/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer(cfg config.HTTP, lifecycle fx.Lifecycle) fiber.Router {
	app := fiber.New(fiber.Config{
		AppName: "Anilibrary Monitoring Server",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(pprof.New())

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.MonitoringPort))

			logging.Get().Info("Monitoring server started at", zap.String("addr", addr))

			go func() {
				if err := app.Listen(addr); err != nil {
					logging.Get().Error("Monitoring server failed to start", zap.Error(err))
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
