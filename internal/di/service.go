package di

import (
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/repository/kafka"
	"anilibrary-scraper/internal/repository/redis"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/service/event"
	scraperService "anilibrary-scraper/internal/service/scraper"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"services",
	fx.Provide(
		fx.Annotate(redis.NewAnimeRepository, fx.As(new(scraperService.AnimeRepository))),
		fx.Annotate(kafka.NewEventRepository, fx.As(new(event.Repository))),

		fx.Annotate(client.NewTLSClient, fx.As(new(scraper.HTTPClient))),
		fx.Annotate(
			func(client scraper.HTTPClient, validator *validator.Validate) scraper.Scraper {
				return scraper.New(
					scraper.WithHTTPClient(client),
					scraper.WithValidator(validator),
					scraper.WithPanicHandler(func() {
						if err := recover(); err != nil {
							metrics.IncrPanicCounter()
						}
					}),
				)
			},
			fx.As(new(scraperService.Scraper)),
		),
	),
)
