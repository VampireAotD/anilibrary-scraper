package router

import (
	"net/http"

	_ "anilibrary-scraper/docs" // generated swagger docs
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/handler/http/router/routes/api"
	"anilibrary-scraper/internal/handler/http/v1/anime"
	"anilibrary-scraper/pkg/logging"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Logger          logging.Contract
	Handler         anime.Controller
	URL             string
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

	api.ComposeRoutes(router, config.Handler)

	return router
}
