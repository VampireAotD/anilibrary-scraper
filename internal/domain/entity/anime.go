package entity

import (
	"encoding/json"
)

type Status string

const (
	Ongoing  Status = "Онгоинг"
	Announce Status = "Анонс"
	Ready    Status = "Вышел"
)

type Anime struct {
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      Status   `json:"status"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voice_acting"`
	Rating      float32  `json:"rating"`
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
