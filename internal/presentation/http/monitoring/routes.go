package monitoring

import (
	"anilibrary-scraper/internal/presentation/http/monitoring/healthcheck"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Router             fiber.Router `name:"monitoring-server"`
	Metrics            *fiberprometheus.FiberPrometheus
	HealthcheckHandler healthcheck.Controller
}

func RegisterMetricsRoutes(params Params) {
	params.Metrics.RegisterAt(params.Router, "/metrics")

	params.Router.Get("/healthcheck", params.HealthcheckHandler.Healthcheck)
}
