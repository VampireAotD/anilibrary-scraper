package router

import (
	"net/http"

	"anilibrary-request-parser/internal/handler/http/api/anime"
	"anilibrary-request-parser/internal/handler/http/routes"
	"anilibrary-request-parser/internal/handler/http/routes/api"
	"github.com/go-chi/chi/v5"
)

func NewRouter(enablePprof bool, controller anime.Controller) http.Handler {
	router := chi.NewRouter()

	router.Handle("/metrics", routes.PrometheusRoutes())

	if enablePprof {
		router.Mount("/debug", routes.ProfilerRoutes())
	}

	api.ComposeRoutes(router, controller)

	return router
}