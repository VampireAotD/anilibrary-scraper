package providers

import (
	"fmt"
	"os"

	"anilibrary-scraper/pkg/logging"
)

const DefaultLoggerFileLocation string = "../../storage/logs/app.log"

func NewLoggerProvider() (logging.Contract, func(), error) {
	file, err := createLogFile()
	if err != nil {
		return nil, nil, err
	}

	logger := logging.NewLogger(os.Stdout, file)
	cleanup := func() {
		logger.Info("closing logger")
		_ = logger.Sync()
		if err := file.Close(); err != nil {
			logger.Error("closing logger file", logging.Error(err))
		}
	}

	logger.Info("Initialized logger")

	return logger, cleanup, nil
}

func createLogFile() (*os.File, error) {
	file, err := os.OpenFile(DefaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, fmt.Errorf("creating log file %w", err)
	}

	return file, nil
}
