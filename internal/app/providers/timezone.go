package providers

import (
	"fmt"
	"time"

	"anilibrary-scraper/pkg/logger"
)

func NewTimezoneProvider(timezone string, log logger.Contract) error {
	log.Info("Setting timezone")

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("setting timezone: %w", err)
	}

	time.Local = location
	return nil
}
