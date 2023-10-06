package providers

import (
	"context"
	"fmt"
	"os"

	"anilibrary-scraper/pkg/logging"

	"go.uber.org/fx"
)

const defaultLoggerFileLocation string = "../../storage/logs/app.log"

func createLogFile() (*os.File, error) {
	file, err := os.OpenFile(defaultLoggerFileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, fmt.Errorf("creating log file %w", err)
	}

	return file, nil
}

func NewLoggerProvider(lifecycle fx.Lifecycle) error {
	file, err := createLogFile()
	if err != nil {
		return fmt.Errorf("creating log file: %w", err)
	}

	logger := logging.New(logging.WithLogFiles(file), logging.ECSCompatible())

	logger.Info("Initialized logger")

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing logger")

			_ = logger.Sync()
			return file.Close()
		},
	})

	return nil
}
