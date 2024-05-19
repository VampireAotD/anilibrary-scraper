package providers

import (
	"context"

	"anilibrary-scraper/pkg/logging"

	"go.uber.org/fx"
)

func NewLoggerProvider(lifecycle fx.Lifecycle) error {
	logger := logging.New(logging.ConvertToJSON(), logging.ECSCompatible(), logging.AsDefault())

	logger.Info("Initialized logger")

	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			logger.Info("Closing logger")

			return logger.Sync()
		},
	})

	return nil
}
