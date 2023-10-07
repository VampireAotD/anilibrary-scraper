package di

import (
	"anilibrary-scraper/config"
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/handler/http/monitoring/healthcheck"
	"anilibrary-scraper/internal/handler/http/router"
	"anilibrary-scraper/internal/handler/http/server"

	"go.uber.org/fx"
)

var HTTPModule = fx.Module(
	"http",
	fx.Provide(
		anime.NewController,
		healthcheck.NewController,
		router.NewHandlers,
		router.New,
		server.New,
	),

	fx.Decorate(func(cfg config.App, router *router.Router) *router.Router {
		if cfg.Env == "local" {
			router.WithProfiling().WithSwagger()
		}

		return router
	}),

	fx.Invoke(func(server server.Server, lifecycle fx.Lifecycle) {
		server.Start(lifecycle)
	}),
)
