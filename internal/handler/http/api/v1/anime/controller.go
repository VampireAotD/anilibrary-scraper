package anime

import (
	"encoding/json"
	"net/http"

	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/pkg/logging"
	"anilibrary-scraper/pkg/response"

	"go.opentelemetry.io/otel/codes"
)

type Controller struct {
	service service.ScraperService
}

func NewController(service service.ScraperService) Controller {
	return Controller{
		service: service,
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
//	@Failure		422				{object}	response.Error
//	@Router			/anime/parse [post]
func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	logger := middleware.MustGetLogger(r.Context())
	tracer := middleware.MustGetTracer(r.Context())
	ctx, span := tracer.Start(r.Context(), "Parse")
	defer span.End()

	resp := response.New(w)

	span.AddEvent("Decoding request")

	var request ScrapeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logger.Error("while decoding incoming request", logging.Error(err))

		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	span.AddEvent("Validating request")

	if err := request.Validate(); err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logger.Error("while decoding incoming url", logging.Error(err))

		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	logger.Info("Scraping", logging.String("url", request.URL))

	span.AddEvent("Scraping data")

	entity, err := c.service.Process(ctx, request.URL)
	if err != nil {
		metrics.IncrHTTPErrorsCounter()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		logger.Error("while scraping", logging.Error(err))

		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	span.AddEvent("Finished scraping")

	metrics.IncrHTTPSuccessCounter()
	_ = resp.JSON(http.StatusOK, entity)
}