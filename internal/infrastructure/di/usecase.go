package di

import (
	"anilibrary-scraper/internal/application/service/event"
	scraperService "anilibrary-scraper/internal/application/service/scraper"
	"anilibrary-scraper/internal/application/usecase/scraper"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Module(
	"usecases",
	fx.Provide(
		fx.Annotate(scraperService.NewScraperService, fx.As(new(scraper.Service))),
		fx.Annotate(event.NewService, fx.As(new(scraper.EventService))),
	),
)
