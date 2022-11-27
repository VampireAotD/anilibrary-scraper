package model

import (
	"encoding/json"

	"anilibrary-scraper/internal/domain/entity"
)

type Anime struct {
	Image       string        `json:"image"`
	Title       string        `json:"title"`
	Status      entity.Status `json:"status"`
	Episodes    string        `json:"episodes"`
	Genres      []string      `json:"genres"`
	VoiceActing []string      `json:"voiceActing"`
	Rating      float32       `json:"rating"`
}

func NewFromEntity(anime *entity.Anime) *Anime {
	return &Anime{
		Image:       anime.Image,
		Title:       anime.Title,
		Status:      anime.Status,
		Episodes:    anime.Episodes,
		Genres:      anime.Genres,
		VoiceActing: anime.VoiceActing,
		Rating:      anime.Rating,
	}
}

func (a *Anime) FromJSON(data []byte) (*Anime, error) {
	err := json.Unmarshal(data, a)

	return a, err
}

func (a *Anime) ToJSON() ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (a *Anime) ToEntity() *entity.Anime {
	return &entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      a.Status,
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Rating:      a.Rating,
	}
}
