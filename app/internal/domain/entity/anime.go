package entity

import (
	"encoding/json"

	"anilibrary-request-parser/app/internal/domain/enum/anime"
)

type Anime struct {
	Title       string       `json:"title"`
	Status      anime.Status `json:"status"`
	Rating      float32      `json:"rating"`
	Episodes    string       `json:"episodes"`
	Genres      []string     `json:"genres"`
	VoiceActing []string     `json:"voice_acting"`
}

func (a Anime) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(a)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
