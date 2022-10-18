package providers

import (
	"fmt"
	"time"

	"anilibrary-scraper/pkg/logging"
)

func NewTimezoneProvider(timezone string, logger logging.Contract) error {
	logger.Info("Setting timezone")

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("setting timezone: %w", err)
	}

	time.Local = location
	return nil
}
