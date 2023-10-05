package di

import (
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/service/event"
	"anilibrary-scraper/internal/domain/service/scraper"
	scraperInstance "anilibrary-scraper/internal/scraper"
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
