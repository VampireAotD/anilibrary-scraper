//go:build wireinject
// +build wireinject

package app

import (
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/v1/anime"
	"anilibrary-scraper/internal/providers"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// Handlers

func WireAnimeController(client *redis.Client) anime.Controller {
	wire.Build(providers.HTTPAnimeHandlerSet)
	return anime.Controller{}
}

// App

func WireDependencies(cfg config.Config) (Dependencies, func(), error) {
	panic(wire.Build(
		providers.NewLoggerProvider,
		wire.FieldsOf(&cfg, "Redis"),
		providers.NewRedisProvider,
		SetupDependencies,
	))
}

func WireApp() (*App, func(), error) {
	panic(wire.Build(
		config.New,
		WireDependencies,
		New,
	))
}
