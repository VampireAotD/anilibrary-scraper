package app

import (
	"os"
	"time"
	_ "time/tzdata"

	"anilibrary-request-parser/pkg/logger"
)

func (a *App) SetTimezone() {
	location, err := time.LoadLocation(a.config.App.Timezone)

	if err != nil {
		defer a.closer.Close()

		a.logger.Error("error while setting timezone", logger.Error(err))

		os.Exit(1)
	}

	time.Local = location
}