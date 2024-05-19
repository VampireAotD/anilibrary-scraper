package di

import (
	"anilibrary-scraper/internal/application/usecase/scraper"
	"anilibrary-scraper/internal/presentation/http/api"
	"anilibrary-scraper/internal/presentation/http/api/v1/anime"
	"anilibrary-scraper/internal/presentation/http/monitoring"
	"anilibrary-scraper/internal/presentation/http/monitoring/healthcheck"

	"github.com/ansrivas/fiberprometheus/v2"
	"go.uber.org/fx"
)

var HTTPModule = fx.Module(
	"http servers",

	fx.Provide(
		fx.Annotate(scraper.NewUseCase, fx.As(new(anime.ScraperUseCase))),
		anime.NewController,
		healthcheck.NewController,
	),

	fx.Supply(fiberprometheus.New("Anilibrary Scraper")),

	fx.Provide(
		fx.Annotate(
			api.NewServer,
			fx.ResultTags(`name:"api-server"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			api.RegisterAPIRoutes,
		),
	),

	fx.Provide(
		fx.Annotate(
			monitoring.NewServer,
			fx.ResultTags(`name:"monitoring-server"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			monitoring.RegisterMetricsRoutes,
		),
	),
)
