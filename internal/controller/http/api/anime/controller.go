package anime

import (
	"anilibrary-request-parser/internal/domain/service/anime"
	"anilibrary-request-parser/pkg/logger"
)

type Controller struct {
	logger  logger.Logger
	service *anime.ScraperService
}

func NewController(logger logger.Logger, service *anime.ScraperService) Controller {
	return Controller{
		logger:  logger,
		service: service,
	}
}