package entity

import (
	"encoding/json"

	"anilibrary-request-parser/internal/domain/enum/anime"
)

type Anime struct {
	Title       string       `json:"title"`
	Status      anime.Status `json:"status"`
	Episodes    string       `json:"episodes"`
	Genres      []string     `json:"genres"`
	VoiceActing []string     `json:"voice_acting"`
	Rating      float32      `json:"rating"`
}

func (a *Anime) FromJSON(data []byte) (*Anime, error) {
	err := json.Unmarshal(data, a)

	return a, err
}

func (a Anime) ToJSON() ([]byte, error) {
	bytes, err := json.Marshal(a)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
