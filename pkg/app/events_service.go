package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
	"github.com/liamg/magic"
)

type EventsService struct{
	eventsRepo 	domain.EventsRepo
	filesRepo 	domain.FilesRepo
}

func NewEventsService(eventsRepo domain.EventsRepo, filesRepo domain.FilesRepo) domain.EventsSvc{
	return &EventsService{
		eventsRepo: eventsRepo,
		filesRepo: filesRepo,
	}
}

func (s *EventsService) CreateEvent(ctx context.Context, eventRequest domain.CreateEventRequest, profilePic *[]byte, creatorId uuid.UUID) error {
	//	create UUID for event
	eventId := uuid.New()

	//	get image extension
	fileType, err := magic.Lookup(*profilePic)
	if err != nil{
		return domain.ErrInvalidImage
	}
	log.Println(fileType.Extension)

	// store profile picutre
	profilePicName := strings.ReplaceAll(eventRequest.Name, " ", "") + "profile"
	path := fmt.Sprintf("events/%s/", eventId)
	url, err  := s.filesRepo.StoreImage(ctx, profilePic, domain.EventsGroup,  profilePicName, fileType.Extension,  path)
	if err != nil{
		return domain.ErrInternal
	}

	//create event
	event := domain.Event{
		Id: eventId,
		Name: eventRequest.Name,
		Description: eventRequest.Description,
		StartsAt: eventRequest.StartsAt,
		EndsAt: eventRequest.EndsAt,
		ProfilePicUrl: url,
		Address: eventRequest.Address,
		CreatedAt: time.Now(),
	}
	if err := s.eventsRepo.CreateEvent(ctx, event); err != nil{
		 _ = s.filesRepo.DeleteFile(ctx, domain.EventsGroup, url)
		return domain.ErrInternal
	}
	return nil
}