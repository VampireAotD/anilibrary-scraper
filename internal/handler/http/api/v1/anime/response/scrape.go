package response

import "anilibrary-scraper/internal/entity"

type Entry struct {
	Name string `json:"name"`
}

type ScrapeResponse struct {
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	Status      string  `json:"status"`
	Episodes    string  `json:"episodes"`
	Genres      []Entry `json:"genres"`
	VoiceActing []Entry `json:"voiceActing"`
	Synonyms    []Entry `json:"synonyms"`
	Rating      float32 `json:"rating"`
}

func NewScrapeResponse(anime entity.Anime) ScrapeResponse {
	return ScrapeResponse{
		Image:       anime.Image,
		Title:       anime.Title,
		Status:      anime.Status,
		Episodes:    anime.Episodes,
		Genres:      mapToEntries(anime.Genres),
		VoiceActing: mapToEntries(anime.VoiceActing),
		Synonyms:    mapToEntries(anime.Synonyms),
		Rating:      anime.Rating,
	}
}

func mapToEntries(data []string) []Entry {
	entries := make([]Entry, 0, len(data))

	for i := range data {
		entries = append(entries, Entry{Name: data[i]})
	}

	return entries
}
