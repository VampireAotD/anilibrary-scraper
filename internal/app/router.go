package app

import (
	"net/http"

	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/internal/routes"
	"anilibrary-request-parser/internal/routes/api"
	"github.com/go-chi/chi/v5"
)

func (app *App) Router() http.Handler {
	router := chi.NewRouter()

	controller, err := app.Controller()

	if err != nil {
		app.stopOnError("composing controller", err)
	}

	api.ComposeRoutes(router, controller)

	router.Handle("/metrics", routes.PrometheusRoutes())

	if app.config.App.Env == config.Local {
		router.Mount("/debug", routes.ProfilerRoutes())
	}

	return router
}
