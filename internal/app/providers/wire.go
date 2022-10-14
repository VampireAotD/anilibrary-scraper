//go:build wireinject
// +build wireinject

package providers

import (
	"anilibrary-scraper/internal/handler/http/api"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
)

// Handlers

func WireAnimeController(client *redis.Client) anime.Controller {
	wire.Build(api.AnimeControllerProviderSet)
	return anime.Controller{}
}
