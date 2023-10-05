package di

import (
	"anilibrary-scraper/internal/domain/usecase"
	"anilibrary-scraper/internal/domain/usecase/scraper"
	"go.uber.org/fx"
)

var UseCaseModule = fx.Module(
	"usecases",
	fx.Provide(
		fx.Annotate(scraper.NewUseCase, fx.As(new(usecase.ScraperUseCase))),
	),
)
