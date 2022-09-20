package app

import (
	"net/http"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/router"
)

func (app *App) Router() http.Handler {
	controller, err := app.Controller()

	if err != nil {
		app.stopOnError("composing controller", err)
	}

	return router.NewRouter(app.config.App.Env == config.Local, controller)
}
