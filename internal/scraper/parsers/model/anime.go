package model

import "anilibrary-scraper/internal/domain/entity"

type Status entity.Status

const (
	Ongoing  Status = "Онгоинг"
	Announce Status = "Анонс"
	Ready    Status = "Вышел"
)

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
