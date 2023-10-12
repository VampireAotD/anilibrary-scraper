package di

import (
	scraperInstance "anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/service"
	"anilibrary-scraper/internal/service/event"
	"anilibrary-scraper/internal/service/scraper"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"services",
	fx.Provide(
		fx.Annotate(scraperInstance.New, fx.As(new(scraperInstance.Contract))),
		fx.Annotate(scraper.NewScraperService, fx.As(new(service.ScraperService))),
		fx.Annotate(event.NewService, fx.As(new(service.EventService))),
	),
)
