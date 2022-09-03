package app

import (
	"os"

	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/pkg/logger"
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
