package app

import (
	"fmt"

	"anilibrary-request-parser/internal/composite"
	"anilibrary-request-parser/internal/handler/http/api/anime"
	"anilibrary-request-parser/internal/providers"
)

func (app *App) Controller() (anime.Controller, error) {
	redisComposite, err := composite.NewRedisComposite(app.config.Redis)

	if err != nil {
		return anime.Controller{}, fmt.Errorf("redis composite: %w", err)
	}

	app.closer.Add("redis composite", redisComposite)

	return providers.WireAnimeController(redisComposite.Client, app.logger.Named("api/http"))
}
