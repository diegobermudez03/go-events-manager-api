package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authSvc 	domain.AuthSvc
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) MountRoutes(router *chi.Mux){
	r := chi.NewRouter()
	r.Post("/register", h.registerUser)
}


//dtos for endpoints
type registerDTO struct{
	fullName	string 	`json:"fullName" validate:"required"`
	age 		int 	`json:"age" validate:"required"`
	gender 		string 	`json:"gender" validate:"required"`
	email 		string 	`json:"email" validate:"required"`
	password 	string 	`json:"password" validate:"required"`
}


func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request){
	if r.Body == nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("No body"))
		return 
	}

	//validate body
	var payload registerDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil{
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid body"))
		return 
	}
	if err := utils.Validate.Struct(payload); err != nil{
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, validationErrors)
		return
	}

	//process request
	refreshToken, err := h.authSvc.RegisterUser(
		payload.fullName,
		payload.age,
		payload.gender,
		payload.email,
		payload.password,
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, map[string]string{
		"refreshToken" : refreshToken,
	})

}