//go:build wireinject
// +build wireinject

package app

import (
	"anilibrary-scraper/internal/app/providers"
	"anilibrary-scraper/internal/config"
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

// App

func WireApp() (*App, func(), error) {
	panic(wire.Build(
		providers.NewLoggerProvider,
		config.New,
		wire.FieldsOf(new(config.Config), "Redis"),
		providers.NewRedisProvider,
		New,
	))
}
