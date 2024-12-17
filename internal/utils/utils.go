package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)


var Validate = validator.New(validator.WithRequiredStructEnabled())

func WriteJSON(w http.ResponseWriter, payload interface{}) {

}

func WriteError(w http.ResponseWriter, status int, err error){

}