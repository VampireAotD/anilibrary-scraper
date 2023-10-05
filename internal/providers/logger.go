package providers

import (
	"context"
	"fmt"
	"os"

	"anilibrary-scraper/pkg/logging"

	"go.uber.org/fx"
)

const DefaultLoggerFileLocation string = "../../storage/logs/app.log"

func NewLoggerProvider(lifecycle fx.Lifecycle) (logging.Contract, error) {
	file, err := createLogFile()
	if err != nil {
		return nil, err
	}

	logger := logging.NewLogger(os.Stdout, file)

	logger.Info("Initialized logger")

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing logger")

			_ = logger.Sync()
			return file.Close()
		},
	})

	return logger, nil
}

func createLogFile() (*os.File, error) {
	file, err := os.OpenFile(DefaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, fmt.Errorf("creating log file %w", err)
	}

	return file, nil
}
