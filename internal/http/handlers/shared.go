package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
)

var (
	ErrNoBody = errors.New("no body")
	ErrInavlidBody = errors.New("invalid body")
)

func validateBody(r *http.Request, payload interface{}) error{
	if r.Body == nil{
		return ErrNoBody
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil{
		return ErrInavlidBody
	}
	if err := utils.Validate.Struct(payload); err != nil{
		//foundErrors := err.(validator.ValidationErrors)
		return ErrInavlidBody
	}
	return nil
}