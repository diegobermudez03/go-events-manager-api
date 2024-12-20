package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/internal/http/middlewares"
	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EventsHandler struct {
	eventsService 	domain.EventsSvc
	middlewares 	*middlewares.Middlewares
}

func NewEventsHandler(eventsService domain.EventsSvc, middlewares 	*middlewares.Middlewares) *EventsHandler {
	return &EventsHandler{
		eventsService: eventsService,
		middlewares: middlewares,
	}
}

func (h *EventsHandler) MountRoutes(router *chi.Mux){
	r := chi.NewRouter()
	r.Use(h.middlewares.AuthMiddleware)
	r.Post("/", h.CreateEventHandler)

	router.Mount("/events", r)
}

/////////		REQUESTS DTOS 
const eventProfileImage = "eventProfile"
const eventBody = "eventBody"

type createEventDTO struct{
	Name 	 	string		`json:"name" validate:"required"`
	Description string		`json:"description" validate:"required"`
	StartsAt 	time.Time 	`json:"startsAt" validate:"required"`
	EndsAt 		time.Time 	`json:"endsAt" validate:"required"`
	Address 	string 		`json:"address" validate:"required"`
}


/////////////////////////////////////////////////////////////////////////
/////////////////////		HANDLERS		////////////////////////////
/////////////////////////////////////////////////////////////////////////


func (h *EventsHandler) CreateEventHandler(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value(middlewares.UserIdKey).(uuid.UUID)

	//Getting the event DTO
	eventBody := r.FormValue(eventBody)
	if eventBody == ""{
		log.Println("No body")
		utils.WriteError(w, http.StatusBadRequest, ErrInavlidBody)
		return 
	}
	var eventDTO createEventDTO
	if err := json.Unmarshal([]byte(eventBody), &eventDTO); err != nil{
		log.Printf("Error with unmarshal %s", err.Error())
		utils.WriteError(w, http.StatusBadRequest, ErrInavlidBody)
		return 
	}

	if err := utils.Validate.Struct(eventDTO); err != nil{
		foundErrors := err.(validator.ValidationErrors)
		log.Printf("error with validation %s", foundErrors)
		utils.WriteError(w, http.StatusBadRequest, ErrInavlidBody)
		return 
	}

	//Getting the profile picture
	err := r.ParseMultipartForm(1024 * 10)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	file, handler, err := r.FormFile(eventProfileImage)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	log.Printf("Filename %s Size %d", handler.Filename, handler.Size)

	//Calling service
	err = h.eventsService.CreateEvent(
		r.Context(),
		domain.CreateEventRequest{
			Name: eventDTO.Name,
			Description: eventDTO.Description,
			StartsAt:  eventDTO.StartsAt,
			EndsAt: eventDTO.EndsAt,
			Address: eventDTO.Address,
		},
		&bytes,
		userId,
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)

}