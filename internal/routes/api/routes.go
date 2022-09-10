package api

import (
	"anilibrary-request-parser/internal/controller/http/api/anime"
	"anilibrary-request-parser/internal/controller/http/middleware"
	services "anilibrary-request-parser/internal/domain/service/anime"
	"anilibrary-request-parser/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func ComposeRoutes(router chi.Router, logger logger.Logger, service *services.ScraperService) {
	controller := anime.NewController(logger.Named("api/http"), service)

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ResponseMetrics)

		r.Route("/anime", func(r chi.Router) {
			r.Post("/parse", controller.Parse)
		})
	})
}
