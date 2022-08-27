package app

import (
	"context"
	"net/http"

	"anilibrary-request-parser/app/internal/composite"
	"anilibrary-request-parser/app/internal/controller/http/v1/anime"
	"anilibrary-request-parser/app/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	redisComposite, err := composite.NewComposite(context.Background(), a.config.Redis)

	if err != nil {
		a.logger.Error("composite", logger.Error(err))
		return nil, err
	}

	a.closer.Add("redisComposite", redisComposite)

	service := composite.NewScraperComposite(redisComposite)
	controller := anime.NewController(a.logger, service)

	composeRoutes(router, a, controller)

	return router, nil
}

func composeRoutes(router chi.Router, app *App, controller anime.Controller) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/anime", func(r chi.Router) {
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
