package anime

import (
	"strings"

	"anilibrary-request-parser/app/internal/domain/contract"
	"anilibrary-request-parser/app/internal/infrastructure/client"
	"anilibrary-request-parser/app/internal/infrastructure/scraper"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animego"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animevost"
	"anilibrary-request-parser/app/pkg/logger"
)

type ScraperService struct {
	scraper contract.Scraper
	logger  logger.Logger
}

func NewScrapperService(url string, logger logger.Logger) (*ScraperService, error) {
	var instance contract.Scraper
	base := scraper.New(url, client.DefaultClient(), logger)

	switch true {
	case strings.Contains(url, "animego.org"):
		instance = animego.New(base)
		break
	case strings.Contains(url, "animevost.org"):
		instance = animevost.New(base)
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
