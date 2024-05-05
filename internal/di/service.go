package di

import (
	"anilibrary-scraper/internal/repository/kafka"
	"anilibrary-scraper/internal/repository/redis"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/service/event"
	scraperService "anilibrary-scraper/internal/service/scraper"

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
