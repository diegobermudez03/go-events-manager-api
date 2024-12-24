package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

/*	having all these interfaces in the same file is not a very good idea for a clean archtieecture
	but for this small project I prefer to keep it like that, its really a overkill to create a lot of
	subpackages only to manage single interfaces on each file, however, in a bigger project I should do something like

	domain/
		users/
			services.go
			repositories.go
			models.go
		orders/
			services.go
			repositories.go
			models.go
		shared/
			custom_errors.go
*/

type AuthSvc interface {
	RegisterUser(ctx context.Context,age int, fullName, gender, email, password string) (string, string, error)
	LoginUser(ctx context.Context, email string, password string) (string, string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	CheckAuthEvent(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, neededPermissions []string) error
}

type UserSvc interface {
	GetUsers(ctx context.Context, filters ...UsersFilter) ([]User, error)
}

type EventsSvc interface{
	CreateEvent(ctx context.Context, event CreateEventRequest, profilePic *[]byte, creatorId uuid.UUID) error 
	GetParticipationsOfUser(ctx context.Context, userId uuid.UUID, filters ...ParticipationFilter) ([]Participation, error)
	GetEvent(ctx context.Context, eventId uuid.UUID)(*EventWithParticipants, error)
	AddParticipation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, roleName string) error
}

type RolesSvc interface{
	GetRoleById(ctx context.Context, roleId uuid.UUID)(*Role, error)
	//GetRoleFromParticipation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) (Role, error)
}

type InitializeSvc interface {
	RegisterRoles() error 
}

// REQUEST TYPES 

//	EventsSvc requests
type CreateEventRequest struct{
	Name 		string 
	Description string 
	StartsAt 	time.Time
	EndsAt		time.Time
	Address 	string 
}
