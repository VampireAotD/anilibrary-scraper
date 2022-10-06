//go:build wireinject
// +build wireinject

package providers

import (
	"anilibrary-scraper/internal/domain/service"
	services "anilibrary-scraper/internal/domain/service/anime"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
)

// Services

func WireScraperService(client *redis.Client) services.ScraperService {
	wire.Build(service.ScraperProviderSet)
	return services.ScraperService{}
}

// Controllers

func WireAnimeController(client *redis.Client, logger logger.Contract) anime.Controller {
	wire.Build(WireScraperService, anime.NewController)
	return anime.Controller{}
}
