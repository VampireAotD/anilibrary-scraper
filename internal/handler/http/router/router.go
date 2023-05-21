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
	"github.com/segmentio/kafka-go"
	swagger "github.com/swaggo/http-swagger"
)

type config struct {
	logger          logging.Contract
	redisConnection *redis.Client
	kafkaConnection *kafka.Conn
	enableProfiling bool
}

func NewRouter(options ...Option) http.Handler {
	cfg := &config{
		enableProfiling: false,
	}

	for i := range options {
		options[i](cfg)
	}

	router := chi.NewRouter()

	router.Use(
		chiMiddleware.Recoverer,
		middleware.Tracer,
		middleware.Logger(cfg.logger),
	)

	router.Handle("/metrics", promhttp.Handler())

	if cfg.enableProfiling {
		router.Mount("/debug", chiMiddleware.Profiler())
	}

	router.Get("/swagger/*", swagger.Handler())

	router.Get("/healthcheck", container.MakeHealthcheckController(cfg.redisConnection, cfg.kafkaConnection).Healthcheck)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ResponseMetrics, middleware.JWTAuth)

		r.Route("/anime", func(r chi.Router) {
			r.Post("/parse", container.MakeAnimeController(cfg.redisConnection, cfg.kafkaConnection).Parse)
		})
	})

	return router
}
