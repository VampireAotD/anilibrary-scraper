package entity

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidData = errors.New("entity was filled with incorrect data")

type Status string

type Anime struct {
	Image       string   `json:"image"`
	Title       string   `json:"title"`
	Status      Status   `json:"status"`
	Episodes    string   `json:"episodes"`
	Genres      []string `json:"genres"`
	VoiceActing []string `json:"voiceActing"`
	Synonyms    []string `json:"synonyms"`
	Rating      float32  `json:"rating"`
}

func (a *Anime) IsValid() error {
	if a.Image == "" || a.Title == "" {
		return ErrInvalidData
	}

	return nil
}

func (a *Anime) FromJSON(data []byte) (*Anime, error) {
	err := json.Unmarshal(data, a)

	return a, err
}

func (a *Anime) ToJSON() ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("marshaling: %w", err)
	}

	return bytes, nil
}
