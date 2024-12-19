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