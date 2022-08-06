package anime

import (
	"encoding/json"
	"errors"
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/utils"
	"anilibrary-request-parser/app/internal/controller/http/v1/anime/dto"
	"anilibrary-request-parser/app/internal/domain/service/anime"
	"anilibrary-request-parser/app/pkg/logger"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	var parse dto.ParseDTO

	json.NewDecoder(r.Body).Decode(&parse)
	defer r.Body.Close()

	err := parse.Validate()

	if err != nil {
		c.logger.Error("while decoding incoming url", logger.Error(err))
		_ = utils.NewError(w, http.StatusUnprocessableEntity, errors.New("invalid url"))
		return
	}

	service, err := anime.NewScraperService(parse.Url)

	if err != nil {
		c.logger.Error("while creating scraper service", logger.Error(err))
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	c.logger.Info("Scraping", logger.String("url", parse.Url))
	entity, err := service.Process()

	if err != nil {
		c.logger.Error("while scraping", logger.Error(err))
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	marshal, _ := json.Marshal(entity)
	w.Write(marshal)
}
