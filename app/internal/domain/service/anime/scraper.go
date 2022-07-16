package anime

import (
	"strings"

	"anilibrary-request-parser/app/internal/infrastructure/scraper"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animego"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animevost"
	"anilibrary-request-parser/app/pkg/logger"
)

type ScraperService struct {
	scraper *scraper.Scrapper
	logger  logger.Logger
}

func NewScrapperService(url string, logger logger.Logger) (*ScraperService, error) {
	var instance *scraper.Scrapper

	switch true {
	case strings.Contains(url, "animego.org"):
		instance = scraper.New(url, animego.New(), logger)
		break
	case strings.Contains(url, "animevost.org"):
		instance = scraper.New(url, animevost.New(), logger)
		break
	default:
		logger.Error("Undefined scraper")
		return nil, nil
	}

	return &ScraperService{
		scraper: instance,
		logger:  logger,
	}, nil
}
