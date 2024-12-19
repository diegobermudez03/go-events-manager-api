package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
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



func domainErrorToHttp(err error) int{
	switch err.Error(){
	case domain.ErrIncorrectPassword.Error(): fallthrough
	case domain.ErrInvalidParametersError.Error(): fallthrough
	case ErrNoBody.Error(): fallthrough
	case ErrInavlidBody.Error(): fallthrough
	case domain.ErrUserDoesntExist.Error(): return http.StatusBadRequest
	
	case domain.ErrInternal.Error(): return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}