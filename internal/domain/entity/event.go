package entity

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	URL  string `json:"url"`
	Date int64  `json:"date"`
}

func (e *Event) ToJSON() ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("marshaling event: %w", err)
	}

	return bytes, nil
}
