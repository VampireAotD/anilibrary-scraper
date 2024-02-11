package model

import (
	"errors"

	"anilibrary-scraper/internal/entity"
)

var ErrInvalidData = errors.New("model doesn't have required data")

type Anime struct {
	URL         string   `json:"url"`
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voiceActing"`
	Synonyms    []string `json:"synonyms"`
	Rating      float32  `json:"rating"`
}

func (a Anime) Validate() error {
	if a.Image == "" || a.Title == "" {
		return ErrInvalidData
	}

	return nil
}

func (a Anime) MapToDomainEntity() entity.Anime {
	return entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      a.Status,
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Synonyms:    a.Synonyms,
		Rating:      a.Rating,
	}
}
