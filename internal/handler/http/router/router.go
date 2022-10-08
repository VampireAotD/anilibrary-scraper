package router

import (
	"net/http"

	"anilibrary-scraper/internal/handler/http/api/anime"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/handler/http/router/routes"
	"anilibrary-scraper/internal/handler/http/router/routes/api"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	Url             string
	EnableProfiling bool
	Logger          logger.Contract
	Handler         anime.Controller
}

func NewRouter(config Config) http.Handler {
	router := chi.NewRouter()

	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.Logger(config.Logger))

	router.Handle("/metrics", routes.PrometheusRoutes())

	if config.EnableProfiling {
		router.Mount("/debug", routes.ProfilerRoutes())
	}

	api.ComposeRoutes(router, config.Handler)

	return router
}
