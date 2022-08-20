package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponse(w http.ResponseWriter, code int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(Error{
		Message: err.Error(),
	})
}
