package app

import (
	"net/http"

	"anilibrary-request-parser/internal/composite"
	"anilibrary-request-parser/internal/controller/http/api/anime"
	"anilibrary-request-parser/internal/routes"
	"anilibrary-request-parser/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func (a *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	redisComposite, err := composite.NewRedisComposite(a.config.Redis)

	if err != nil {
		a.logger.Error("redis composite", logger.Error(err))
		return nil, err
	}

	a.closer.Add("redis composite", redisComposite)

	service := composite.NewScraperComposite(redisComposite)
	controller := anime.NewController(a.logger.Named("api/http"), service)

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
		router.Handle("/metrics", routes.PrometheusRoutes())
		app.logger.Info("Prometheus metrics enabled", logger.String("endpoint", "/metrics"))
	}

	if app.flags.pprof {
		router.Mount("/debug", routes.ProfilerRoutes())
		app.logger.Info("Pprof enabled", logger.String("endpoint", "/debug/pprof"))
	}
}
