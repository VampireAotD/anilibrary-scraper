package router

import (
	"net/http"

	_ "anilibrary-scraper/docs" // generated swagger docs
	"anilibrary-scraper/internal/container"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/pkg/logging"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Logger          logging.Contract
	RedisConnection *redis.Client
	EnableProfiling bool
}

func NewRouter(config *Config) http.Handler {
	router := chi.NewRouter()

	router.Use(
		chiMiddleware.Recoverer,
		middleware.Tracer,
		middleware.Logger(config.Logger),
	)

	router.Handle("/metrics", promhttp.Handler())

	if config.EnableProfiling {
		router.Mount("/debug", chiMiddleware.Profiler())
	}

	router.Get("/swagger/*", swagger.Handler())

	router.Get("/healthcheck", container.MakeHealthcheckController(config.RedisConnection).Healthcheck)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ResponseMetrics, middleware.JWTAuth)

		r.Route("/anime", func(r chi.Router) {
			r.Post("/parse", container.MakeAnimeController(config.RedisConnection).Parse)
		})
	})

	return router
}
