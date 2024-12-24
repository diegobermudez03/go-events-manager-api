package handlers

import (
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/http/middlewares"
	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	usersService domain.UserSvc
	middelwares  *middlewares.Middlewares
}

// REQUESTS URL QUERIES
const text = "text"
const limit = "limit"
const offset = "offset"

func NewUsersHandler(usersService domain.UserSvc, middelwares  *middlewares.Middlewares) *UsersHandler{
	return &UsersHandler{
		usersService: usersService,
		middelwares: middelwares,
	}
}


func (h *UsersHandler) MountRoutes(router *chi.Mux){
	r := chi.NewRouter()
	r.Use(h.middelwares.AuthMiddleware)
	
	r.Get("/", h.GetUsersByEmailOrName)

	router.Mount("/users", r)
}


func (h *UsersHandler) GetUsersByEmailOrName(w http.ResponseWriter, r *http.Request){
	textFilter := r.URL.Query().Get(text)
	var textPointer *string
	if textFilter != "" {
		textPointer = &textFilter 
	}

	limitFilter := getIntQueryParam(limit, r)
	offsetFilter := getIntQueryParam(offset, r)

	users, err := h.usersService.GetUsers(
		r.Context(), 
		domain.UsersTextFilter(textPointer),
		domain.UsersOffsetFilter(offsetFilter),
		domain.UsersLimitFilter(limitFilter),
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}

	utils.WriteJSON(w, http.StatusOK, users)
}