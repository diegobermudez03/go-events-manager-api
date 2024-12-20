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
	r.Get("/", h.GetEventsFromUser)

	router.Mount("/events", r)
}

/////////		REQUESTS DTOS 
const eventProfileImage = "eventProfile"
const eventBody = "eventBody"

const eventOffsetQuery = "offset"
const eventLimitQuery = "limit"
const eventRoleQuery = "role"

type createEventDTO struct{
	Name 	 	string		`json:"name" validate:"required"`
	Description string		`json:"description" validate:"required"`
	StartsAt 	time.Time 	`json:"startsAt" validate:"required"`
	EndsAt 		time.Time 	`json:"endsAt" validate:"required"`
	Address 	string 		`json:"address" validate:"required"`
}

//////////	RESPONSE DTOS
type eventTileDTO struct{
	Id 				uuid.UUID	`json:"id"`
	Name 			string		`json:"name"`
	Description 	string 		`json:"description"`
	StartsAt		time.Time	`json:"startsAt"`
	RoleName 		string 		`json:"roleName"`
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

func (h *EventsHandler) GetEventsFromUser(w http.ResponseWriter, r *http.Request){
	userId, ok := r.Context().Value(middlewares.UserIdKey).(uuid.UUID)
	if !ok{
		utils.WriteError(w, http.StatusBadRequest, ErrInavlidBody)
		return 
	}
	//Getting query params
	limit := getIntQueryParam(eventLimitQuery, r)
	offset := getIntQueryParam(eventOffsetQuery, r)
	roleAux := r.URL.Query().Get(eventRoleQuery)
	role := new(string)
	if roleAux != ""{
		*role = roleAux
	}else{
		role = nil
	}

	//call service
	participations, err := h.eventsService.GetParticipationsOfUser(
		r.Context(),
		userId,
		domain.ParticipationRoleFilter(role),
		domain.ParticipationOffsetFilter(offset),
		domain.ParticipationLimitFilter(limit),
	)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	eventDTOs := make([]eventTileDTO, len(participations))
	for index, ev := range participations{
		eventDTOs[index].Id = ev.Event.Id
		eventDTOs[index].Name = ev.Event.Name
		eventDTOs[index].Description = ev.Event.Description
		eventDTOs[index].StartsAt = ev.Event.StartsAt
		eventDTOs[index].RoleName = ev.RoleName
	}
	utils.WriteJSON(w, http.StatusOK, eventDTOs)
}
