package anime

import (
	"encoding/json"
	"net/http"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/pkg/logging"
	"anilibrary-scraper/pkg/response"
)

type Controller struct {
	service service.ScraperService
}

func NewController(service service.ScraperService) Controller {
	return Controller{
		service: service,
	}
}

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	log := middleware.MustGetLogger(r.Context())
	tracer := middleware.MustGetTracer(r.Context())
	_, span := tracer.Start(r.Context(), "Parse")
	defer span.End()

	resp := response.New(w)
	parseDTO := dto.RequestDTO{
		FromCache: true,
	}

	json.NewDecoder(r.Body).Decode(&parseDTO)
	if err := parseDTO.Validate(); err != nil {
		span.RecordError(err)
		log.Error("while decoding incoming url", logging.Error(err))

		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	log.Info("Scraping", logging.String("url", parseDTO.Url))

	entity, err := c.service.Process(parseDTO)
	if err != nil {
		span.RecordError(err)
		log.Error("while scraping", logging.Error(err))

		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	_ = resp.JSON(http.StatusOK, entity)
}
