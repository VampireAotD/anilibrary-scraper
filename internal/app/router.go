package app

import (
	"net/http"

	"anilibrary-request-parser/internal/routes"
	"anilibrary-request-parser/internal/routes/api"
	"anilibrary-request-parser/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func (app *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	animeService, err := app.AnimeService()

	if err != nil {
		return nil, err
	}

	api.ComposeRoutes(router, app.logger, animeService)

	if app.flags.prom {
		router.Handle("/metrics", routes.PrometheusRoutes())
		app.logger.Info("Prometheus metrics enabled", logger.String("endpoint", "/metrics"))
	}

	if app.flags.pprof {
		router.Mount("/debug", routes.ProfilerRoutes())
		app.logger.Info("Pprof enabled", logger.String("endpoint", "/debug/pprof"))
	}

	return router, nil
}
