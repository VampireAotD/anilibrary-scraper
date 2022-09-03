package app

import (
	"log"

	"anilibrary-request-parser/pkg/logger"
)

func (a *App) SetLogger() {
	instance, err := logger.New(a.flags.logPath)

	if err != nil {
		log.Fatal(err)
	}

	a.logger = instance.Logger
	a.closer.Add("logger", instance)
}
