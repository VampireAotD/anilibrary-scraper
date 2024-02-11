package anime

import (
	"context"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/pkg/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//go:generate mockgen -source=controller.go -destination=./mocks.go -package=anime
type ScraperUseCase interface {
	Scrape(ctx context.Context, url string) (*entity.Anime, error)
}

type Controller struct {
	useCase ScraperUseCase
}

func NewController(useCase ScraperUseCase) Controller {
	return Controller{
		useCase: useCase,
	}
}

// Parse
//
//	@Summary		Scrape anime data
//	@Description	Scrape anime data
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Access token"	default(Bearer)
//	@Param			url				body		ScrapeRequest	true	"Url to scrape from"
//	@Success		200				{object}	ScrapeResponse
//	@Failure		401				string		Unauthorized
//	@Failure		422				{object}	ErrorResponse
//	@Router			/anime/parse [post]
func (c Controller) Parse(ctx *fiber.Ctx) error {
	span := trace.SpanFromContext(ctx.UserContext())
	defer span.End()

	span.AddEvent("Decoding and validating request")

	var request ScrapeRequest
	if err := request.MapAndValidate(ctx); err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(NewErrorResponse(err))
	}

	logging.Get().Info("Scraping", zap.String("url", request.URL))
	span.AddEvent("Scraping data")

	anime, err := c.useCase.Scrape(ctx.UserContext(), request.URL)
	if err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logging.Get().Error("while scraping", zap.Error(err))

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(NewErrorResponse(err))
	}

	span.AddEvent("Finished scraping")

	metrics.IncrHTTPSuccessCounter()
	return ctx.JSON(EntityToResponse(anime))
}
