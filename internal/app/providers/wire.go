//go:build wireinject
// +build wireinject

package providers

import (
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
)

// Services

func WireScraperService(client *redis.Client) scraper.Service {
	wire.Build(service.ScraperProviderSet)
	return scraper.Service{}
}

// Controllers

func WireAnimeController(client *redis.Client, logger logger.Contract) anime.Controller {
	wire.Build(WireScraperService, anime.NewController)
	return anime.Controller{}
}
