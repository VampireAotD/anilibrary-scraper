package api

import (
	"anilibrary-request-parser/internal/handler/http/api/anime"
	"anilibrary-request-parser/internal/handler/http/middleware"
	"github.com/go-chi/chi/v5"
)

func ComposeRoutes(router chi.Router, controller anime.Controller) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ResponseMetrics)

		r.Route("/anime", func(r chi.Router) {
			r.Post("/parse", controller.Parse)
		})
	})
}
