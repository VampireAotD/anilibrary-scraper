package app

import (
	"fmt"

	"anilibrary-request-parser/internal/composite"
	"anilibrary-request-parser/internal/domain/service/anime"
)

func (app *App) AnimeService() (*anime.ScraperService, error) {
	redisComposite, err := composite.NewRedisComposite(app.config.Redis)

	if err != nil {
		return nil, fmt.Errorf("redis composite: %w", err)
	}

	app.closer.Add("redis composite", redisComposite)

	return composite.NewScraperComposite(redisComposite), nil
}
