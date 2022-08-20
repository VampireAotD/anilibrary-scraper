package app

import (
	"os"

	"anilibrary-request-parser/app/internal/config"
	"anilibrary-request-parser/app/pkg/logger"
)

func (a *App) ReadConfig() {
	cfg, err := config.New()

	if err != nil {
		defer a.closer.Close()

		a.logger.Error("error while reading config", logger.Error(err))

		os.Exit(1)
	}

	a.config = cfg
}
