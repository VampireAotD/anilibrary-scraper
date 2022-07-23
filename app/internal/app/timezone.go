package app

import (
	"os"
	"time"

	"go.uber.org/zap"
)

func (a *App) SetTimezone() {
	location, err := time.LoadLocation(a.config.App.Timezone)

	if err != nil {
		defer a.closer.Close()

		a.logger.Error("error while setting timezone", zap.Error(err))

		os.Exit(1)
	}

	time.Local = location
}
