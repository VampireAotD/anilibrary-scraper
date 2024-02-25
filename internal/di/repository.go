package di

import (
	"anilibrary-scraper/internal/repository/kafka"
	"anilibrary-scraper/internal/repository/redis"
	"anilibrary-scraper/internal/service/event"
	"anilibrary-scraper/internal/service/scraper"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Module(
	"repositories",
	fx.Provide(
		fx.Annotate(redis.NewAnimeRepository, fx.As(new(scraper.AnimeRepository))),
		fx.Annotate(kafka.NewEventRepository, fx.As(new(event.Repository))),
	),
)
