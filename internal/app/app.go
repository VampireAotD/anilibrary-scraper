package app

import (
	"os"
	"time"

	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/pkg/closer"
	"anilibrary-request-parser/pkg/logger"
)

type App struct {
	logger logger.Logger
	config config.Config
	closer closer.Closers
}

func (app *App) stopOnError(info string, err error) {
	app.logger.Error(info, logger.Error(err))
	os.Exit(1)
}

func (app *App) ReadConfig() {
	cfg, err := config.New()

	if err != nil {
		app.stopOnError("error while reading config", err)
	}

	app.config = cfg
}

func (app *App) SetTimezone() {
	location, err := time.LoadLocation(app.config.App.Timezone)

	if err != nil {
		app.stopOnError("error while setting timezone", err)
	}

	time.Local = location
}
