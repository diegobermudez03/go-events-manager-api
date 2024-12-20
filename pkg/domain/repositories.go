package domain

import (
	"context"

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

type AuthRepo interface {
	GetUserAuthByEmail(ctx context.Context, email string) (*UserAuth, error)
	GetUserAuthById(ctx context.Context, id uuid.UUID) (*UserAuth, error)
	RegisterUser(ctx context.Context, user UserAuth) error
}

type UsersRepo interface{
	CreateUser(ctx context.Context, user User) error
}

type SessionsRepo interface {
	CreateSession(ctx context.Context, session Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error
}

type RolesRepo interface{
	CreateRoleIfNotExists(ctx context.Context, role Role) error
	GetRoleByName(ctx context.Context, roleName string) (*Role, error)
	GetRoleIdByName(ctx context.Context, roleName string) (uuid.UUID, error)
}

type EventsRepo interface{
	CreateEvent(ctx context.Context, event Event) error
	CreateParticipant(ctx context.Context, userId uuid.UUID, eventId uuid.UUID, roleId uuid.UUID) error
	GetParticipations(ctx context.Context, filters ParticipationFilters) ([]Participation, error)
}

type FilesRepo interface{
	// GROUP IS IN CASE WE HAVE AN IMPLEMENTATION LIKE S3 or MINIO, IT WOULD REPRESENT THE BUCKET
	StoreImage(ctx context.Context, image *[]byte, group string, imageName string, imageType string, path string) (string, error)
	DeleteFile(ctx context.Context, group string,url string) error
}