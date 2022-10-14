package anime

import (
	"encoding/json"
	"net/http"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/pkg/logger"
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
	resp := response.New(w)
	log := middleware.MustGetLogger(r.Context())
	parseDTO := dto.RequestDTO{
		FromCache: true,
	}

	json.NewDecoder(r.Body).Decode(&parseDTO)
	if err := parseDTO.Validate(); err != nil {
		log.Error("while decoding incoming url", logger.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	log.Info("Scraping", logger.String("url", parseDTO.Url))

	entity, err := c.service.Process(parseDTO)
	if err != nil {
		log.Error("while scraping", logger.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, err)
		return
	}

	_ = resp.JSON(http.StatusOK, entity)
}
