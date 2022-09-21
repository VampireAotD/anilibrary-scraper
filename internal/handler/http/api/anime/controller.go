package anime

import (
	"anilibrary-scraper/internal/domain/service/anime"
	"anilibrary-scraper/pkg/logger"
)

type Controller struct {
	logger  logger.Contract
	service *anime.ScraperService
}

func NewController(logger logger.Contract, service *anime.ScraperService) Controller {
	return Controller{
		logger:  logger,
		service: service,
	}
}
