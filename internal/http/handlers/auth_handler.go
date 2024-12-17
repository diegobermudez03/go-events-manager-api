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

func NewAuthHandler(authSvc domain.AuthSvc) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

func (h *AuthHandler) MountRoutes(router *chi.Mux){
	r := chi.NewRouter()
	r.Post("/register", h.registerUser)

	router.Mount("/auth", r)
}


//dtos for endpoints
type registerDTO struct{
	FullName	string 	`json:"fullName" validate:"required"`
	Age 		int 	`json:"age" validate:"required"`
	Gender 		string 	`json:"gender" validate:"required"`
	Email 		string 	`json:"email" validate:"required,email"`
	Password 	string 	`json:"password" validate:"required"`
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
	refreshToken, accessToken, err := h.authSvc.RegisterUser(
		r.Context(),
		payload.Age,
		payload.FullName,
		payload.Gender,
		payload.Email,
		payload.Password,
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(
		w, 
		http.StatusCreated,
		map[string]string{
			"refreshToken" : refreshToken,
			"accessToken" : accessToken,
		},
	)
}