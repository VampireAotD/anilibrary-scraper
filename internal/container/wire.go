//go:build wireinject
// +build wireinject

package container

import (
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/handler/http/monitoring/healthcheck"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func MakeAnimeController(client *redis.Client, kafka *kafka.Conn) anime.Controller {
	panic(wire.Build(HTTPAnimeHandlerSet))
}

func MakeHealthcheckController(client *redis.Client, kafka *kafka.Conn) healthcheck.Controller {
	panic(wire.Build(healthcheck.NewController))
}
