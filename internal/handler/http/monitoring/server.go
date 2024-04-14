package monitoring

import (
	"context"
	"net"
	"strconv"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/logging"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

func NewServer(lifecycle fx.Lifecycle, cfg config.HTTP) fiber.Router {
	app := fiber.New(fiber.Config{
		AppName: "Anilibrary Monitoring Server",
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logging.Get(),
	}))
	app.Use(recover.New())
	app.Use(pprof.New())

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.MonitoringPort))

			logging.Get().Info("Monitoring server started at", logging.String("addr", addr))

			go func() {
				if err := app.Listen(addr); err != nil {
					logging.Get().Error("Monitoring server failed to start", logging.Error(err))
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
