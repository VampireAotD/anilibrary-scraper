package monitoring

import (
	"github.com/VampireAotD/anilibrary-scraper/internal/presentation/http/monitoring/healthcheck"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	HealthcheckHandler healthcheck.Controller
	Router             fiber.Router `name:"monitoring-server"`
	Metrics            *fiberprometheus.FiberPrometheus
}

func RegisterMetricsRoutes(params Params) {
	params.Metrics.RegisterAt(params.Router, "/metrics")

	params.Router.Get("/healthcheck", params.HealthcheckHandler.Healthcheck)
}
