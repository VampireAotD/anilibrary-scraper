package anime

import (
	"encoding/json"
	"errors"
	"net/http"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/pkg/logger"
	"anilibrary-scraper/pkg/response"
)

type Controller struct {
	service scraper.Service
}

func NewController(service scraper.Service) Controller {
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
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	defer r.Body.Close()

	log.Info("Scraping", logger.String("url", parseDTO.Url))

	entity, err := c.service.Process(parseDTO)
	if err != nil {
		log.Error("while scraping", logger.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	_ = resp.JSON(http.StatusOK, entity)
}
