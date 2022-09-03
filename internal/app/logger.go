package app

import (
	"log"
	"os"

	"anilibrary-request-parser/pkg/logger"
)

const DefaultLoggerFileLocation string = "../../storage/logs/app.log"

func (a *App) SetLogger() {
	file := createLogFile()

	instance, err := logger.New(logger.Config{
		ConsoleOutput: os.Stdout,
		LogFile:       file,
	})

	if err != nil {
		log.Fatal(err)
	}

	a.logger = instance
	a.closer.Add("logger", file)
}

func createLogFile() *os.File {
	file, err := os.OpenFile(DefaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatalf("error while creating file %s", err)
	}

	return file
}
