package app

import (
	"anilibrary-request-parser/pkg/closer"
)

func (a *App) SetCloser() {
	a.closer = closer.New(a.logger)
}
