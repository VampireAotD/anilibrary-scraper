package composite

import (
	"anilibrary-request-parser/internal/domain/repository/redis"
	services "anilibrary-request-parser/internal/domain/service/anime"
)

type ScraperComposite struct {
	*services.ScraperService
}

func NewScraperComposite(composite RedisComposite) *services.ScraperService {
	animeRepository := redis.NewAnimeRepository(composite.client)

	return services.NewScraperService(animeRepository)
}
