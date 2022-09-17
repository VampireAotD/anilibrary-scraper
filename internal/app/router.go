package app

import (
	"net/http"

	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/internal/handler/http/router"
)

func (app *App) Router() http.Handler {
	controller, err := app.Controller()

	if err != nil {
		app.stopOnError("composing controller", err)
	}

	return router.NewRouter(app.config.App.Env == config.Local, controller)
}
