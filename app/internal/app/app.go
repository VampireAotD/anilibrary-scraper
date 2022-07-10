package app

import (
	"anilibrary-request-parser/app/internal/config"
	"anilibrary-request-parser/app/pkg/closer"
	"anilibrary-request-parser/app/pkg/logger"
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
