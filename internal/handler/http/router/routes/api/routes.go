package api

import (
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/handler/http/v1/anime"
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
