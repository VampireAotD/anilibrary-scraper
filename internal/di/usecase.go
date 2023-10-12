package di

import (
	"anilibrary-scraper/internal/usecase"
	"anilibrary-scraper/internal/usecase/scraper"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Module(
	"usecases",
	fx.Provide(
		fx.Annotate(scraper.NewUseCase, fx.As(new(usecase.ScraperUseCase))),
	),
)
