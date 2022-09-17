package anime

import (
	"encoding/json"
	"errors"
	"net/http"

	"anilibrary-request-parser/internal/domain/dto"
	"anilibrary-request-parser/pkg/logger"
	"anilibrary-request-parser/pkg/response"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	resp := response.New(w)

	parseDTO := dto.ParseDTO{
		FromCache: true,
	}

	json.NewDecoder(r.Body).Decode(&parseDTO)
	err := parseDTO.Validate()

	if err != nil {
		c.logger.Error("while decoding incoming url", logger.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	defer r.Body.Close()

	c.logger.Info("Scraping", logger.String("url", parseDTO.Url))
	entity, err := c.service.Process(parseDTO)

	if err != nil {
		c.logger.Error("while scraping", logger.Error(err))
		_ = resp.ErrorJSON(http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	_ = resp.JSON(http.StatusOK, entity)
}
