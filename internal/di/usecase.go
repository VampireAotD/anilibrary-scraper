package di

import (
	"anilibrary-scraper/internal/service/event"
	scraperService "anilibrary-scraper/internal/service/scraper"
	"anilibrary-scraper/internal/usecase/scraper"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Module(
	"usecases",
	fx.Provide(
		fx.Annotate(scraperService.NewScraperService, fx.As(new(scraper.Service))),
		fx.Annotate(event.NewService, fx.As(new(scraper.EventService))),
	),
)
