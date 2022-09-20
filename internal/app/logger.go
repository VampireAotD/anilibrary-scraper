package app

import (
	"log"
	"os"

	"anilibrary-request-parser/pkg/logger"
)

const DefaultLoggerFileLocation string = "../../storage/logs/app.log"

func (app *App) InitLogger() {
	file := createLogFile()

	app.logger = logger.NewLogger(os.Stdout, file)
	app.closer.Add("logger", func() error {
		_ = app.logger.Sync()
		return file.Close()
	})
}

func createLogFile() *os.File {
	file, err := os.OpenFile(DefaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatalf("error while creating file %s", err)
	}

	return file
}
