package utils

import (
	"encoding/json"
	"net/http"

	"github.com/creasty/defaults"
)

func FromJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON[T any](w http.ResponseWriter, status int, v *T) {
	w.WriteHeader(status)
	logger := GetLogger("[utils-json] ")

	if v == nil {
		logger.Println("WriteJSON called with nil response body")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := defaults.Set(v); err != nil {
		logger.Printf("Could not apply defaults: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Printf("Could not write json response: %v\n", err)
	}
}
