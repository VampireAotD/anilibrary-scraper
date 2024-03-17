package anime

import (
	"context"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/handler/http/api/v1/anime/request"
	"anilibrary-scraper/internal/handler/http/api/v1/anime/response"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/usecase/scraper"
	"anilibrary-scraper/pkg/logging"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//go:generate mockgen -source=controller.go -destination=./mocks.go -package=anime
type ScraperUseCase interface {
	Scrape(ctx context.Context, dto scraper.DTO) (entity.Anime, error)
}

type Controller struct {
	useCase   ScraperUseCase
	validator *validator.Validate
}

func NewController(useCase ScraperUseCase, validate *validator.Validate) Controller {
	return Controller{
		useCase:   useCase,
		validator: validate,
	}
}

// Parse
//
//	@Summary		Scrape anime data
//	@Description	Scrape anime data
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Access token"	default(Bearer)
//	@Param			url				body		request.ScrapeRequest	true	"Url to scrape from"
//	@Success		200				{object}	response.ScrapeResponse
//	@Failure		401				string		Unauthorized
//	@Failure		422				{object}	response.ErrorResponse
//	@Router			/anime/parse [post]
func (c Controller) Parse(ctx *fiber.Ctx) error {
	span := trace.SpanFromContext(ctx.UserContext())
	defer span.End()

	span.AddEvent("Decoding and validating request")

	var req request.ScrapeRequest
	if err := req.MapAndValidate(ctx, c.validator); err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewErrorResponse(err))
	}

	logging.Get().Info("Scraping", zap.String("url", req.URL))
	span.AddEvent("Scraping data")

	anime, err := c.useCase.Scrape(ctx.UserContext(), scraper.DTO{
		URL:       req.URL,
		IP:        ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	})
	if err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logging.Get().Error("while scraping", zap.Error(err))

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewErrorResponse(err))
	}

	span.AddEvent("Finished scraping")

	metrics.IncrHTTPSuccessCounter()
	return ctx.JSON(response.NewScrapeResponse(anime))
}
