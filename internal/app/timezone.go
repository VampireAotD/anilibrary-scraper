package app

import (
	"os"
	"time"
	_ "time/tzdata"

	"anilibrary-request-parser/pkg/logger"
)

func (app *App) SetTimezone() {
	location, err := time.LoadLocation(app.config.App.Timezone)

	if err != nil {
		defer app.closer.Close()

		app.logger.Error("error while setting timezone", logger.Error(err))

		os.Exit(1)
	}

	time.Local = location
}
