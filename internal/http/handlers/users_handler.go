package handlers

import "github.com/go-chi/chi/v5"

type UsersHandler struct {
}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}

func (h *UsersHandler) MountRoutes(router *chi.Mux){
	
}