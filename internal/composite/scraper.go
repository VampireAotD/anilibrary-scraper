package composite

import (
	"anilibrary-request-parser/internal/adapter/db/redis/repository"
	services "anilibrary-request-parser/internal/domain/service/anime"
)

type ScraperComposite struct {
	*services.ScraperService
}

func NewScraperComposite(composite RedisComposite) *services.ScraperService {
	animeRepository := repository.NewAnimeRepository(composite.client)

	return services.NewScraperService(animeRepository)
}
