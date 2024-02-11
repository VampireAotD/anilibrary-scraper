package anime

import "anilibrary-scraper/internal/entity"

type ScrapeResponse struct {
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voiceActing"`
	Synonyms    []string `json:"synonyms"`
	Rating      float32  `json:"rating"`
}

func EntityToResponse(anime entity.Anime) ScrapeResponse {
	return ScrapeResponse(anime)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}
