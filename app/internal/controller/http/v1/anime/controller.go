package anime

import "anilibrary-request-parser/app/pkg/logger"

type Controller struct {
	logger logger.Logger
}

func NewController(logger logger.Logger) *Controller {
	return &Controller{logger: logger}
}
