package anime

import (
	"context"
	"errors"

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

// Scrape
//
//	@Summary		Scrape anime data
//	@Description	Scrape anime data
//	@Tags			anime
//	@Accept			json
//	@Produce		json
//	@Param			url	body		request.ScrapeRequest	true	"Url to scrape from"
//	@Success		200	{object}	response.ScrapeResponse
//	@Failure		401	string		Unauthorized
//	@Failure		422	{object}	response.ScrapeErrorResponse
//	@Router			/anime/scrape [post]
//	@Security		Bearer
func (c Controller) Scrape(ctx *fiber.Ctx) error {
	logger := logging.FromContext(ctx.UserContext())
	span := trace.SpanFromContext(ctx.UserContext())
	defer span.End()

	span.AddEvent("Parsing and validating request")

	var req request.ScrapeRequest
	if err := req.Validate(ctx, c.validator); err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, "failed to parse HTTP request")
		span.RecordError(err)

		if errors.Is(err, request.ErrUnableToDecodeRequest) {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewScrapeError("Invalid request type"))
		}

		if errors.Is(err, request.ErrInvalidURL) {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewScrapeError("Invalid URL"))
		}

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewScrapeError(err.Error()))
	}

	logger.Info("Scraping", logging.String("url", req.URL))
	span.AddEvent("Scraping anime")

	anime, err := c.useCase.Scrape(ctx.UserContext(), scraper.DTO{
		URL:       req.URL,
		IP:        ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	})
	if err != nil {
		metrics.IncrHTTPErrorsCounter()
		logger.Error("Failed to scrape anime", logging.Error(err))
		span.SetStatus(codes.Error, "failed to scrape anime from HTTP request")
		span.RecordError(err)

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response.NewScrapeError(err.Error()))
	}

	span.AddEvent("Finished scraping")

	metrics.IncrHTTPSuccessCounter()
	return ctx.JSON(response.NewScrapeResponse(anime))
}
