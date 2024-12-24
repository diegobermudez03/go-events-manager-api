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
	authRepo 	domain.AuthRepo
	emailService domain.EmailSvc
	invitationsBus domain.InvitationsEventBus
}

func NewEventsService(
	eventsRepo domain.EventsRepo, 
	rolesRepo 	domain.RolesRepo,
	usersRepo 	domain.UsersRepo,
	filesRepo domain.FilesRepo,
	authRepo 	domain.AuthRepo,
	emailService domain.EmailSvc,
	invitationsBus domain.InvitationsEventBus) domain.EventsSvc{
	return &EventsService{
		eventsRepo: eventsRepo,
		rolesRepo: rolesRepo,
		usersRepo: usersRepo,
		filesRepo: filesRepo,
		authRepo: authRepo,
		emailService: emailService,
		invitationsBus: invitationsBus,
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

	//if no limit passed, then one added
	if filter.Limit == nil{
		filter.Limit = new(int)
		*filter.Limit = 100
	}

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
			roleData, err := s.rolesRepo.GetRoleById(ctx, dataModel.RoleId)
			if err != nil{
				return nil, domain.ErrInternal
			}
			rolesMap[roleData.Id] = roleData.Name
			role = roleData.Name
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

func (s *EventsService) GetEvent(ctx context.Context, eventId uuid.UUID)(*domain.EventWithParticipants, error){
	event, err := s.eventsRepo.GetEventById(ctx, eventId)
	if err != nil{
		return nil, err 
	}

	partFilters := domain.ParticipationFilters{}

	//adding eventId filter
	domain.ParticipationEventIdFilter(&eventId)(&partFilters)

	dataParticipations, err := s.eventsRepo.GetParticipations(ctx, partFilters)
	if err != nil{
		return nil, domain.ErrInternal
	}

	//get participants
	participants := make([]domain.Participant, len(dataParticipations))

	rolesCache := map[uuid.UUID]string{}

	for index, part := range dataParticipations{
		var roleName string 
		if role, ok := rolesCache[part.RoleId]; ok{
			roleName = role 
		}else{
			dtRole, err := s.rolesRepo.GetRoleById(ctx, part.RoleId)
			if err != nil{
				return nil, domain.ErrInternal
			}
			roleName = dtRole.Name
			rolesCache[part.RoleId] = dtRole.Name
		}
		
		user, err := s.usersRepo.GetUserById(ctx, part.UserId)
		if err != nil{
			return nil, domain.ErrInternal
		}

		participants[index] = domain.Participant{
			User: *user,
			Role: roleName,
		}
	}

	eventWithParticipants := domain.EventWithParticipants{
		Event: *event,
		Participants: participants,
	}
	return &eventWithParticipants, nil
}


func (s *EventsService) AddParticipation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, roleName string) error{
	role, err := s.rolesRepo.GetRoleByName(ctx, roleName)
	if err != nil{
		return domain.ErrRoleDoesntExist
	}

	_, err = s.usersRepo.GetUserById(ctx, userId)
	if err != nil{
		return domain.ErrUserDoesntExist
	}

	_, err = s.eventsRepo.GetEventById(ctx, eventId)
	if err != nil{
		return domain.ErrEventDoesntExist
	}

	if err := s.eventsRepo.CreateParticipant(ctx, userId, eventId, role.Id); err != nil{
		return err 
	}
	return nil
}

func (s *EventsService) InviteUser(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error {
	// First check if the invitation already exists, in which case, omit everything
	if exists, _ := s.eventsRepo.CheckInvitation(ctx, eventId, userId); exists{
		return domain.ErrAlreadyInvited
	}
	event, err := s.eventsRepo.GetEventById(ctx, eventId); 
	if err != nil{
		return err 
	}

	userAuth, err := s.authRepo.GetUserAuthById(ctx, userId)
	if err != nil{
		return err 
	}

	if err := s.eventsRepo.CreateInvitation(ctx, eventId, userId); err != nil{
		return err
	}

	//send email in background, we hope it reaches the receiver xd, wont wait for confirmation
	go s.emailService.SendTextEmail(
		ctx, 
		userAuth.Email,
		fmt.Sprintf("You have been invited to %s", event.Name),
		fmt.Sprintf("You have been invited to event %s which will take place on %s at %v", event.Name, event.Address, event.StartsAt),
	)

	s.invitationsBus.Publish(eventId, domain.InvitationEvent{
		Email: userAuth.Email,
		Time: time.Now(),
	})
	return nil
}