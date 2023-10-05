package di

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/kafka"
	"anilibrary-scraper/internal/domain/repository/redis"
	"go.uber.org/fx"
)

var RepositoryModule = fx.Module(
	"repositories",
	fx.Provide(
		fx.Annotate(redis.NewAnimeRepository, fx.As(new(repository.AnimeRepository))),
		fx.Annotate(kafka.NewEventRepository, fx.As(new(repository.EventRepository))),
	),
)
