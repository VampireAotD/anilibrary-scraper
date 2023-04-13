//go:build wireinject
// +build wireinject

package container

import (
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/handler/http/monitoring/healthcheck"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func MakeAnimeController(client *redis.Client) anime.Controller {
	panic(wire.Build(HTTPAnimeHandlerSet))
}

func MakeHealthcheckController(client *redis.Client) healthcheck.Controller {
	panic(wire.Build(healthcheck.NewController))
}
