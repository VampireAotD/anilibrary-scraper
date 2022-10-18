package providers

import (
	"fmt"
	"os"

	"anilibrary-scraper/pkg/logger"
)

const DefaultLoggerFileLocation string = "../../storage/logs/app.log"

func NewLoggerProvider() (logger.Contract, func(), error) {
	file, err := createLogFile()
	if err != nil {
		return nil, nil, err
	}

	logs := logger.NewLogger(os.Stdout, file)
	cleanup := func() {
		logs.Info("closing logger")

		_ = logs.Sync()
		if err := file.Close(); err != nil {
			logs.Error("closing logger file", logger.Error(err))
		}
	}

	logs.Info("Initialized logger")

	return logs, cleanup, nil
}

func createLogFile() (*os.File, error) {
	file, err := os.OpenFile(DefaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, fmt.Errorf("creating log file %w", err)
	}

	return file, nil
}
