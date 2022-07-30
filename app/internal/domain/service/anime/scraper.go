package anime

import (
	"errors"
	"strings"

	"anilibrary-request-parser/app/internal/domain/contract"
	"anilibrary-request-parser/app/internal/infrastructure/client"
	"anilibrary-request-parser/app/internal/infrastructure/scraper"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animego"
	"anilibrary-request-parser/app/internal/infrastructure/scraper/animevost"
)

type ScraperService struct {
	scraper contract.Scraper
}

func NewScrapperService(url string) (*ScraperService, error) {
	var instance contract.Scraper
	base := scraper.New(url, client.DefaultClient())

	switch true {
	case strings.Contains(url, "animego.org"):
		instance = animego.New(base)
		break
	case strings.Contains(url, "animevost.org"):
		instance = animevost.New(base)
		break
	default:
		return nil, errors.New("undefined scraper")
	}

	return &ScraperService{
		scraper: instance,
	}, nil
}
