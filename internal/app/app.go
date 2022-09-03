package app

import (
	"anilibrary-request-parser/internal/config"
	"anilibrary-request-parser/pkg/closer"
	"anilibrary-request-parser/pkg/logger"
)

type App struct {
	flags  flags
	logger logger.Logger
	config *config.Config
	closer closer.Closers
}

func (a *App) Run() {
	a.Listen()
}
