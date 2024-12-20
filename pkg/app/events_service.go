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
	rolesRepo 	domain.RolesRepo
	filesRepo 	domain.FilesRepo
}

func NewEventsService(
	eventsRepo domain.EventsRepo, 
	rolesRepo 	domain.RolesRepo,
	filesRepo domain.FilesRepo) domain.EventsSvc{
	return &EventsService{
		eventsRepo: eventsRepo,
		rolesRepo: rolesRepo,
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
		return err
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
		return err
	}
	//	get creator role and create participant with creator
	roleId, err := s.rolesRepo.GetRoleIdByName(ctx, domain.RoleCreator)	// CONTESTANT FOR CACHE
	if err != nil{
		return err 
	}
	if err := s.eventsRepo.CreateParticipant(ctx, creatorId, eventId, roleId); err != nil{
		return err
	}
	return nil
}

func (s *EventsService) GetParticipationsOfUser(ctx context.Context, userId uuid.UUID, filters ...domain.ParticipationFilter) ([]domain.Participation, error){
	filter := domain.ParticipationFilters{}
	for _, f := range filters{
		f(&filter)
	}
	//explicetely adding userID filter
	domain.ParticipationUserIdFilter(&userId)

	participations, err := s.eventsRepo.GetParticipations(ctx, filter)
	if err != nil{
		return nil, err 
	}
	return participations, nil
}