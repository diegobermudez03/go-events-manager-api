package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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


func getIntQueryParam(key string, r *http.Request) (*int){
	if val := r.URL.Query().Get(key); val != ""{
		num, err := strconv.Atoi(val)
		if err != nil{
			return nil 
		}else{
			p := new(int)
			*p = num
			return p
		}
	}
	return nil
}