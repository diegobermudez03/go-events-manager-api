package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)


var Validate = validator.New(validator.WithRequiredStructEnabled())

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) error{
	return WriteJSON(
		w, 
		status,
		map[string]string{
		"error" : err.Error(),
	})
}
