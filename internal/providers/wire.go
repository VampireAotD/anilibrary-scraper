//go:build wireinject
// +build wireinject

package providers

import (
	"anilibrary-request-parser/internal/controller/http/api/anime"
	"anilibrary-request-parser/internal/domain/service"
	services "anilibrary-request-parser/internal/domain/service/anime"
	"anilibrary-request-parser/pkg/logger"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
)

// Services

func WireScraperService(client *redis.Client) (*services.ScraperService, error) {
	wire.Build(service.ScraperProviderSet)
	return &services.ScraperService{}, nil
}

// Controllers

func WireAnimeController(client *redis.Client, logger logger.Logger) (anime.Controller, error) {
	wire.Build(WireScraperService, anime.NewController)
	return anime.Controller{}, nil
}
