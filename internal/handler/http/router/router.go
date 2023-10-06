package router

import (
	"net/http"

	_ "anilibrary-scraper/docs" // generated swagger docs
	"anilibrary-scraper/internal/handler/http/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Router struct {
	engine   chi.Router
	handlers Handlers
}

func New(handlers Handlers) *Router {
	router := chi.NewRouter()

	router.Use(
		chiMiddleware.Recoverer,
		middleware.Tracer,
	)

	return &Router{
		engine:   router,
		handlers: handlers,
	}
}

func (r *Router) WithMetrics() *Router {
	r.engine.Handle("/metrics", promhttp.Handler())

	return r
}

func (r *Router) WithProfiling() *Router {
	r.engine.Mount("/debug", chiMiddleware.Profiler())

	return r
}

func (r *Router) WithSwagger() *Router {
	r.engine.Get("/swagger/*", swagger.Handler())

	return r
}

func (r *Router) WithLogger(logger *zap.Logger) *Router {
	r.engine.Use(middleware.Logger(logger.Named("http")))

	return r
}

func (r *Router) Routes() http.Handler {
	r.engine.Get("/healthcheck", r.handlers.healthcheck.Healthcheck)

	r.engine.Route("/api/v1", func(router chi.Router) {
		router.Use(middleware.ResponseMetrics, middleware.JWTAuth)

		router.Route("/anime", func(router chi.Router) {
			router.Post("/parse", r.handlers.anime.Parse)
		})
	})

	return r.engine
}
