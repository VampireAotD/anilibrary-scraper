package model

import (
	"errors"

	"anilibrary-scraper/internal/domain/entity"
)

var ErrInvalidParsedData = errors.New("invalid parsed data")

const (
	Ongoing  Status = "Онгоинг"
	Announce Status = "Анонс"
	Ready    Status = "Вышел"
)

type Status entity.Status

type Anime struct {
	Image       string
	Title       string
	Status      Status
	Episodes    string
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Rating      float32
}

func (a *Anime) IsValid() bool {
	return a.Image != "" && a.Title != ""
}

func (a *Anime) ToEntity() *entity.Anime {
	return &entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      entity.Status(a.Status),
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Synonyms:    a.Synonyms,
		Rating:      a.Rating,
	}
}
