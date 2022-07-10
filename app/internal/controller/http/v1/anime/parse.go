package anime

import (
	"encoding/json"
	"net/http"

	"anilibrary-request-parser/app/internal/controller/http/utils"
	"anilibrary-request-parser/app/internal/controller/http/v1/anime/dto"
	"anilibrary-request-parser/app/internal/infrastructure/scrapper"
)

func (c Controller) Parse(w http.ResponseWriter, r *http.Request) {
	var parse dto.ParseDTO

	err := json.NewDecoder(r.Body).Decode(&parse)

	if err != nil {
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	defer r.Body.Close()

	s, err := scrapper.New(parse.Url)

	if err != nil {
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	anime, err := s.Process()

	if err != nil {
		_ = utils.NewError(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(anime.Title))
}
