package handlers

import (
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/go-chi/chi/v5"
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
	r.Get("/login", h.LoginUser)

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

type loginDTO struct{
	Email		string `json:"email" validate:"required"`
	Password 	string `json:"passord" validate:"required"`
}

type loginResponseDTO struct{
	RefreshToken	string	`json:"refreshToken"`
	AccessToken 	string 	`json:"accessToken"`
}



func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request){
	var payload registerDTO
	err := validateBody(r, &payload)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
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
		loginResponseDTO{
			RefreshToken: refreshToken,
			AccessToken: accessToken,
		},
	)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request){
	var payload loginDTO
	err := validateBody(r, &payload)
	if err != nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return 
	}
	refreshToken, accessToken, err := h.authSvc.LoginUser(r.Context(), payload.Email, payload.Password)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(
		w,
		http.StatusAccepted,
		loginResponseDTO{
			RefreshToken: refreshToken,
			AccessToken: accessToken,
		},
	)

	
}