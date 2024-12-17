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

type UsersRepo interface {
	GetUserAuthByEmail(ctx context.Context, email string) (*UserAuth, error)
	RegisterUser(ctx context.Context, auth UserAuth, user User) (uuid.UUID, error) 
}

type SessionsRepo interface {
	CreateSession(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error 
}