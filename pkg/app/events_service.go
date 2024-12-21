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
	usersRepo 	domain.UsersRepo
	rolesRepo 	domain.RolesRepo
	filesRepo 	domain.FilesRepo
}

func NewEventsService(
	eventsRepo domain.EventsRepo, 
	rolesRepo 	domain.RolesRepo,
	usersRepo 	domain.UsersRepo,
	filesRepo domain.FilesRepo) domain.EventsSvc{
	return &EventsService{
		eventsRepo: eventsRepo,
		rolesRepo: rolesRepo,
		usersRepo: usersRepo,
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

	partDataModels, err := s.eventsRepo.GetParticipations(ctx, filter)
	if err != nil{
		return nil, err 
	}
	usersMap := map[uuid.UUID]*domain.User{}
	eventsMap := map[uuid.UUID]*domain.Event{}
	rolesMap := map[uuid.UUID]string{}
	participations := make([]domain.Participation, len(partDataModels))

	//iterate over datamodels and construct entitites
	for index, dataModel := range partDataModels{
		var user *domain.User
		var event *domain.Event
		var role  string
		//get user
		if userAux, ok := usersMap[dataModel.UserId]; ok{
			user = userAux
		}else{
			user, err = s.usersRepo.GetUserById(ctx, dataModel.UserId)
			if err != nil{
				return nil, domain.ErrInternal
			}
			usersMap[user.Id] = user 
		}
		//get event
		if eventAux, ok := eventsMap[dataModel.EventId]; ok{
			event = eventAux
		}else{
			event, err = s.eventsRepo.GetEventById(ctx, dataModel.EventId)
			if err != nil{
				return nil, domain.ErrInternal
			}
			eventsMap[event.Id] = event
		}
		//get role
		if name, ok := rolesMap[dataModel.RoleId]; ok{
			role = name
		}else{
			roleEntity, err := s.rolesRepo.GetRoleById(ctx, dataModel.RoleId)
			if err != nil{
				return nil, domain.ErrInternal
			}
			rolesMap[roleEntity.Id] = roleEntity.Name
		}

		participations[index] = domain.Participation{
			Id : dataModel.Id,
			Event: event,
			User: user,
			RoleName: role,
		}
	}
	return participations, nil
}