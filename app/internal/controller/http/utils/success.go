package utils

import (
	"encoding/json"
	"net/http"
)

func NewSuccessResponse(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(data)
}
