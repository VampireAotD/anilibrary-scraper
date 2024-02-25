package di

import (
	scraperInstance "anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/service/event"
	"anilibrary-scraper/internal/service/scraper"
	scraperUseCase "anilibrary-scraper/internal/usecase/scraper"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"services",
	fx.Provide(
		fx.Annotate(scraperInstance.New, fx.As(new(scraper.Scraper))),
		fx.Annotate(scraper.NewScraperService, fx.As(new(scraperUseCase.Service))),
		fx.Annotate(event.NewService, fx.As(new(scraperUseCase.EventService))),
	),
)
