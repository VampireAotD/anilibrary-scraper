package app

import "anilibrary-request-parser/app/pkg/closer"

func (a *App) SetCloser() {
	a.closer = closer.New(a.logger)
}
