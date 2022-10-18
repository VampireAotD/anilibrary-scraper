package router

import (
	"net/http"

	"anilibrary-scraper/internal/handler/http/api/anime"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/handler/http/router/routes/api"
	"anilibrary-scraper/pkg/logging"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Url             string
	EnableProfiling bool
	Logger          logging.Contract
	Handler         anime.Controller
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

	api.ComposeRoutes(router, config.Handler)

	return router
}
