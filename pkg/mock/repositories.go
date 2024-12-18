package mock

import (
	"context"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

//	AUTH REPO
type AuthRepoMock struct{}
func (u *AuthRepoMock) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error) {
	if email == "d1@gmail.com" {
		return &domain.UserAuth{}, nil
	}
	return nil, domain.ErrUserWithEmailAlreadyExists
}
func (u *AuthRepoMock) RegisterUser(ctx context.Context, auth domain.UserAuth) error {
	return nil
}

//	USERS REPO
type UsersRepoMock struct{}
func (u *UsersRepoMock) CreateUser(ctx context.Context, user domain.User) error {
	return nil
}


//	SESSIONS REPO
type SessionRepoMock struct{}

func (s SessionRepoMock) CreateSession(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error {
	return nil
}