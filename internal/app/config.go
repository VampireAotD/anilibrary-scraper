package app

import (
	"os"

	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/pkg/logger"
)

func (app *App) ReadConfig() {
	cfg, err := config.New()

	if err != nil {
		defer app.closer.Close()

		app.logger.Error("error while reading config", logger.Error(err))

		os.Exit(1)
	}

	app.config = cfg
}
