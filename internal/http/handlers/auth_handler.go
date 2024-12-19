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
	r.Get("/login", h.loginUser)
	r.Get("/refresh", h.refreshSession)

	router.Mount("/auth", r)
}


/////////		REQUESTS DTOS 
type registerDTO struct{
	FullName	string 	`json:"fullName" validate:"required"`
	Age 		int 	`json:"age" validate:"required"`
	Gender 		string 	`json:"gender" validate:"required"`
	Email 		string 	`json:"email" validate:"required,email"`
	Password 	string 	`json:"password" validate:"required"`
}

type loginDTO struct{
	Email		string `json:"email" validate:"required"`
	Password 	string `json:"password" validate:"required"`
}

type refreshDTO struct{
	RefreshToken 	string 	`json:"refreshToken" validate:"required"`
}

/////////		RESPONSES DTOS 
type loginResponseDTO struct{
	RefreshToken	string	`json:"refreshToken"`
	AccessToken 	string 	`json:"accessToken"`
}

type refreshResponseDTO struct{
	AccessToken	string 	`json:"accessToken"`
}

/////////////////////////////////////////////////////////////////////////
/////////////////////		HANDLERS		////////////////////////////
/////////////////////////////////////////////////////////////////////////

func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request){
	var payload registerDTO
	err := validateBody(r, &payload)
	if err != nil{
		utils.WriteError(w, domainErrorToHttp(err), err)
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
		utils.WriteError(w, domainErrorToHttp(err), err)
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

func (h *AuthHandler) loginUser(w http.ResponseWriter, r *http.Request){
	var payload loginDTO
	err := validateBody(r, &payload)
	if err != nil{
		utils.WriteError(w, domainErrorToHttp(err), err)
		return 
	}
	refreshToken, accessToken, err := h.authSvc.LoginUser(r.Context(), payload.Email, payload.Password)
	if err != nil{
		utils.WriteError(w, domainErrorToHttp(err), err)
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

func (h *AuthHandler) refreshSession(w http.ResponseWriter, r *http.Request){
	var payload refreshDTO
	if err := validateBody(r, &payload); err != nil{
		utils.WriteError(w, domainErrorToHttp(err), err)
		return
	}
	accessToken, err:= h.authSvc.RefreshAccessToken(r.Context(), payload.RefreshToken)

	if err != nil{
		utils.WriteError(w, domainErrorToHttp(err), err)
		return 
	}
	utils.WriteJSON(
		w, 
		http.StatusCreated,
		refreshResponseDTO{
			AccessToken: accessToken,
		},
	)
}