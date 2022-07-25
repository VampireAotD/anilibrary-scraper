package app

import (
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/v1/anime"
	"anilibrary-request-parser/app/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	composeRoutes(router, a)

	return router, nil
}

func composeRoutes(router chi.Router, app *App) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/anime", func(r chi.Router) {
			controller := anime.NewController(app.logger)

			r.Post("/parse", controller.Parse)
		})
	})

	if app.flags.prom {
		router.Handle("/metrics", promhttp.Handler())

		app.logger.Info("Prometheus metrics enabled", logger.String("endpoint", "/metrics"))
	}

	if app.flags.pprof {
		router.Mount("/debug", middleware.Profiler())
		app.logger.Info("Pprof enabled", logger.String("endpoint", "/debug/pprof"))
	}
}
