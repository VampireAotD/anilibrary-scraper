package app

import (
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/v1/anime"
	"github.com/go-chi/chi/v5"
)

func (a *App) Router() (http.Handler, error) {
	router := chi.NewRouter()

	composeRoutes(router)

	return router, nil
}

func composeRoutes(router chi.Router) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/anime", func(r chi.Router) {
			controller := anime.NewController()

			r.Post("/parse", controller.Parse)
		})
	})
}
