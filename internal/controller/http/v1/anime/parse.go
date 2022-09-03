package anime

import (
	"encoding/json"
	"errors"
	"net/http"

	utils2 "anilibrary-request-parser/internal/controller/http/utils"
	"anilibrary-request-parser/internal/domain/dto"
	"anilibrary-request-parser/pkg/logger"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	var parseDTO dto.ParseDTO
	parseDTO.FromCache = true

	json.NewDecoder(r.Body).Decode(&parseDTO)
	err := parseDTO.Validate()

	if err != nil {
		c.logger.Error("while decoding incoming url", logger.Error(err))
		_ = utils2.NewErrorResponse(w, http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	defer r.Body.Close()

	c.logger.Info("Scraping", logger.String("url", parseDTO.Url))
	entity, err := c.service.Process(parseDTO)

	if err != nil {
		c.logger.Error("while scraping", logger.Error(err))
		_ = utils2.NewErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	_ = utils2.NewSuccessResponse(w, entity)
}
