package model

import (
	"anilibrary-scraper/internal/domain/entity"
)

type Status = entity.Status
type Type = entity.Type

type Anime struct {
	URL         string   `json:"url"`
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      Status   `json:"status"`
	Type        Type     `json:"type"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voiceActing"`
	Synonyms    []string `json:"synonyms"`
	Rating      float32  `json:"rating"`
	Year        int      `json:"year"`
}

func (a *Anime) MapToDomainEntity() entity.Anime {
	return entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      a.Status,
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Synonyms:    a.Synonyms,
		Rating:      a.Rating,
		Year:        a.Year,
		Type:        a.Type,
	}
}
