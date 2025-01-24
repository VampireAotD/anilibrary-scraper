package response

import (
	"github.com/VampireAotD/anilibrary-scraper/internal/domain/entity"
)

type Entry struct {
	Name string `json:"name"`
}

type ScrapeResponse struct {
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	Genres      []Entry `json:"genres"`
	VoiceActing []Entry `json:"voiceActing"`
	Synonyms    []Entry `json:"synonyms"`
	Episodes    int     `json:"episodes"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
}

func NewScrapeResponse(anime entity.Anime) ScrapeResponse {
	return ScrapeResponse{
		Image:       anime.Image,
		Title:       anime.Title,
		Status:      string(anime.Status),
		Episodes:    anime.Episodes,
		Genres:      mapToEntries(anime.Genres),
		VoiceActing: mapToEntries(anime.VoiceActing),
		Synonyms:    mapToEntries(anime.Synonyms),
		Rating:      anime.Rating,
		Year:        anime.Year,
		Type:        string(anime.Type),
	}
}

func mapToEntries(data []string) []Entry {
	entries := make([]Entry, 0, len(data))

	for i := range data {
		entries = append(entries, Entry{Name: data[i]})
	}

	return entries
}

type ScrapeErrorResponse struct {
	Message string `json:"message"`
}

func NewScrapeError(message string) ScrapeErrorResponse {
	return ScrapeErrorResponse{
		Message: message,
	}
}
