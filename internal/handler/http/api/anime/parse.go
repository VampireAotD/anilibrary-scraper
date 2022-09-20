package anime

import (
	"encoding/json"
	"errors"
	"net/http"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/pkg/response"
	"go.uber.org/zap"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	resp := response.New(w)

	parseDTO := dto.ParseDTO{
		FromCache: true,
	}

	json.NewDecoder(r.Body).Decode(&parseDTO)
	err := parseDTO.Validate()

	if err != nil {
		c.logger.Error("while decoding incoming url", zap.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	defer r.Body.Close()

	c.logger.Info("Scraping", zap.String("url", parseDTO.Url))
	entity, err := c.service.Process(parseDTO)

	if err != nil {
		c.logger.Error("while scraping", zap.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	_ = resp.JSON(http.StatusOK, entity)
}
