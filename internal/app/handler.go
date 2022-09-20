package app

import (
	"fmt"

	"anilibrary-scraper/internal/handler/http/api/anime"
	"anilibrary-scraper/internal/providers"
)

func (app *App) Controller() (anime.Controller, error) {
	redisProvider, err := providers.NewRedisProvider(app.config.Redis)

	if err != nil {
		return anime.Controller{}, fmt.Errorf("redis client: %w", err)
	}

	app.closer.Add("redis client", redisProvider.Close)

	return providers.WireAnimeController(redisProvider, app.logger.Named("api/http"))
}
