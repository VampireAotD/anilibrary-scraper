package di

import (
	"anilibrary-scraper/internal/repository"
	"anilibrary-scraper/internal/repository/kafka"
	"anilibrary-scraper/internal/repository/redis"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Module(
	"repositories",
	fx.Provide(
		fx.Annotate(redis.NewAnimeRepository, fx.As(new(repository.AnimeRepository))),
		fx.Annotate(kafka.NewEventRepository, fx.As(new(repository.EventRepository))),
	),
)
