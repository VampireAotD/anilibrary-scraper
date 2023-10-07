package anime

import (
	"net/http"

	"anilibrary-scraper/internal/domain/usecase"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/metrics"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	usecase usecase.ScraperUseCase
}

func NewController(usecase usecase.ScraperUseCase) Controller {
	return Controller{
		usecase: usecase,
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
//	@Success		200				{object}	entity.Anime
//	@Failure		401				string		Unauthorized
//	@Failure		422				{object}	ErrorResponse
//	@Router			/anime/parse [post]
func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	var (
		logger = middleware.MustGetLogger(r.Context())
		tracer = middleware.MustGetTracer(r.Context())
	)

	ctx, span := tracer.Start(r.Context(), "Parse")
	defer span.End()

	span.AddEvent("Decoding and validating request")

	var request ScrapeRequest
	if err := request.MapAndValidate(r.Body); err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logger.Error("while decoding incoming url", zap.Error(err))

		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, NewErrorResponse(err))
		return
	}

	logger.Info("Scraping", zap.String("url", request.URL))
	span.AddEvent("Scraping data")

	entity, err := c.usecase.Scrape(ctx, request.URL)
	if err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logger.Error("while scraping", zap.Error(err))

		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, NewErrorResponse(err))
		return
	}

	span.AddEvent("Finished scraping")

	metrics.IncrHTTPSuccessCounter()
	render.JSON(w, r, entity)
	return
}
