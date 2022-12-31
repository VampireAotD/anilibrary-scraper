//go:build wireinject
// +build wireinject

package app

import (
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/v1/anime"
	"anilibrary-scraper/internal/providers"

	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
)

// Handlers

func WireAnimeController(client *redis.Client) anime.Controller {
	wire.Build(providers.HTTPAnimeHandlerSet)
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
