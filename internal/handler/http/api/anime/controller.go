package anime

import (
	"anilibrary-request-parser/internal/domain/service/anime"
	"go.uber.org/zap"
)

type Controller struct {
	logger  *zap.Logger
	service *anime.ScraperService
}

func NewController(logger *zap.Logger, service *anime.ScraperService) Controller {
	return Controller{
		logger:  logger,
		service: service,
	}
}
