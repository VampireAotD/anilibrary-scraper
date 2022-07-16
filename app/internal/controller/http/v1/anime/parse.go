package anime

import (
	"encoding/json"
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/utils"
	"anilibrary-request-parser/app/internal/controller/http/v1/anime/dto"
	"anilibrary-request-parser/app/internal/domain/service/anime"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	var parse dto.ParseDTO

	err := json.NewDecoder(r.Body).Decode(&parse)

	if err != nil {
		c.logger.Error("while decoding incoming url")
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	defer r.Body.Close()

	service, err := anime.NewScrapperService(parse.Url, c.logger)

	if err != nil {
		c.logger.Error("while creating scraper service")
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	entity, err := service.Process()

	if err != nil {
		c.logger.Error("while scraping")
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	marshal, _ := json.Marshal(entity)
	w.Write(marshal)
}
