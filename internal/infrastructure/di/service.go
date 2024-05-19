package di

import (
	"anilibrary-scraper/internal/application/service/event"
	scraperService "anilibrary-scraper/internal/application/service/scraper"
	"anilibrary-scraper/internal/infrastructure/repository/kafka"
	"anilibrary-scraper/internal/infrastructure/repository/redis"
	"anilibrary-scraper/internal/infrastructure/scraper"
	"anilibrary-scraper/internal/infrastructure/scraper/client"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"services",
	fx.Provide(
		fx.Annotate(redis.NewAnimeRepository, fx.As(new(scraperService.AnimeRepository))),
		fx.Annotate(kafka.NewEventRepository, fx.As(new(event.Repository))),

		fx.Annotate(client.NewTLSClient, fx.As(new(scraper.HTTPClient))),
		fx.Annotate(scraper.New, fx.As(new(scraperService.Scraper))),
	),
)
