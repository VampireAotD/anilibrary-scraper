package app

import (
	"os"

	"anilibrary-request-parser/app/internal/config"
	"go.uber.org/zap"
)

func (a *App) ReadConfig() {
	cfg, err := config.New(a.flags.envPath)

	if err != nil {
		defer a.closer.Close()

		a.logger.Error("error while reading config", zap.Error(err))

		os.Exit(1)
	}

	a.config = cfg
}
