//go:build wireinject
// +build wireinject

package providers

import (
	"anilibrary-scraper/internal/domain/service"
	services "anilibrary-scraper/internal/domain/service/anime"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// Services

func WireScraperService(client *redis.Client) (*services.ScraperService, error) {
	wire.Build(service.ScraperProviderSet)
	return &services.ScraperService{}, nil
}

// Controllers

func WireAnimeController(client *redis.Client, logger *zap.Logger) (anime.Controller, error) {
	wire.Build(WireScraperService, anime.NewController)
	return anime.Controller{}, nil
}
