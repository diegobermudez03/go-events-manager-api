package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	invitationsBus 	domain.InvitationsEventBus
	middlewares 	*middlewares.Middlewares
}

func NewEventsHandler(eventsService domain.EventsSvc, invitationsBus domain.InvitationsEventBus, middlewares *middlewares.Middlewares) *EventsHandler {
	return &EventsHandler{
		eventsService: eventsService,
		middlewares: middlewares,
		invitationsBus: invitationsBus,
	}
}

func (h *EventsHandler) MountRoutes(router *chi.Mux){
	r := chi.NewRouter()
	r.Use(h.middlewares.AuthMiddleware)
	r.Post("/", h.CreateEventHandler)
	r.Get("/", h.GetEventsFromUser)
	
	r.With(h.middlewares.EventAccessMiddleware(domain.PermissionEditEvent)).
		Get(fmt.Sprintf("/{%s}", eventId), h.GetEvent)

	r.With(h.middlewares.EventAccessMiddleware(domain.PermissionAddParticipant)). 
		Post(fmt.Sprintf("/{%s}/participants", eventId), h.PostParticipant)

	r.With(h.middlewares.EventAccessMiddleware(domain.PermissionInvitePeople)).
		Post(fmt.Sprintf("/{%s}/invitations", eventId), h.PostInvitation)

	r.With(h.middlewares.EventAccessMiddleware(domain.PermissionInvitePeople)).
		Get(fmt.Sprintf("/{%s}/invitations/live", eventId), h.LiveInvitationsToEvent)

	router.Mount("/events", r)
}

/////////		REQUESTS DTOS 
const eventProfileImage = "eventProfile"
const eventBody = "eventBody"

const eventOffsetQuery = "offset"
const eventLimitQuery = "limit"
const eventRoleQuery = "role"

const eventId = "eventId"

type createEventDTO struct{
	Name 	 	string		`json:"name" validate:"required"`
	Description string		`json:"description" validate:"required"`
	StartsAt 	int64 		`json:"startsAt" validate:"required"`
	EndsAt 		int64 		`json:"endsAt" validate:"required"`
	Address 	string 		`json:"address" validate:"required"`
}

type addParticipantDTO struct{
	UserId 		uuid.UUID	`json:"userId" validate:"required"`
	Role 		string 		`json:"role" validate:"required"`
}

//////////	RESPONSE DTOS
type eventTileDTO struct{
	Id 				uuid.UUID	`json:"id"`
	Name 			string		`json:"name"`
	Description 	string 		`json:"description"`
	StartsAt		int64		`json:"startsAt"`
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
			StartsAt:  time.Time(time.Unix(eventDTO.StartsAt, 0)),
			EndsAt:  time.Time(time.Unix(eventDTO.EndsAt, 0)),
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
		eventDTOs[index].StartsAt = ev.Event.StartsAt.Unix()
		eventDTOs[index].RoleName = ev.RoleName
	}
	utils.WriteJSON(w, http.StatusOK, eventDTOs)
}


func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request){
	reqEventId, ok := r.Context().Value(eventId).(uuid.UUID)
	if !ok{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	event, err := h.eventsService.GetEvent(r.Context(), reqEventId)
	if err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *EventsHandler) PostParticipant(w http.ResponseWriter, r *http.Request){
	reqEventId, ok := r.Context().Value(eventId).(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return 
	}

	// extract payload
	if r.Body == nil{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("no body"))
		return 
	}
	var payload addParticipantDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("invalid body"))
		return 
	}
	if err := utils.Validate.Struct(payload); err != nil{
		//foundErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusInternalServerError, errors.New("invalid body"))
		return 
	}

	//call service
	if err := h.eventsService.AddParticipation(r.Context(), reqEventId, payload.UserId, payload.Role); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *EventsHandler) PostInvitation(w http.ResponseWriter, r *http.Request){
	reqEventId, ok := r.Context().Value(eventId).(uuid.UUID)
	if !ok{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return 
	}

	if r.Body == nil{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("no body"))
		return 
	}

	payload := struct{
		UserId 	uuid.UUID	`json:"userId" validate:"required"`
	}{}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("invalid body"))
		return 
	}

	if err := h.eventsService.InviteUser(r.Context(), reqEventId, payload.UserId); err != nil{
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}

	utils.WriteJSON(w, http.StatusAccepted, nil )
}

func (h *EventsHandler) LiveInvitationsToEvent(w http.ResponseWriter, r *http.Request){
	reqEventId, ok := r.Context().Value(eventId).(uuid.UUID)
	if !ok{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return 
	}

	flusher, ok := w.(http.Flusher)
	if !ok{
		utils.WriteError(w, http.StatusInternalServerError, errors.New("streaming unsuported"))
		return 
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	incomingChannel := h.invitationsBus.Suscribe(reqEventId)

	loop := true
	for loop{
		select{
		case event:= <- incomingChannel:
			bytes, err := json.Marshal(event)
			if err != nil{
				continue 
			}
			jsonText := string(bytes)
			w.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonText)))
			flusher.Flush()
		case <-r.Context().Done():
			h.invitationsBus.Unsuscribe(reqEventId, incomingChannel)
			loop = false
		}
	}
}