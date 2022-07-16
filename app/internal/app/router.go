package app

import (
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/v1/anime"
	"anilibrary-request-parser/app/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func (a *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	composeRoutes(router, a.logger)

	return router, nil
}

func composeRoutes(router chi.Router, logger logger.Logger) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/anime", func(r chi.Router) {
			controller := anime.NewController(logger)

			r.Post("/parse", controller.Parse)
		})
	})
}
