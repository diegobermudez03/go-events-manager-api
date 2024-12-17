package app

import (
	"context"
	"testing"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

func TestAuthService(t *testing.T) {

	userMock := UserRepoMock{}
	sessionmock := SessionRepoMock{}

	service := NewAuthService(&userMock, &sessionmock, 1440, 600, "secret")

	t.Run("should receive both tokens", func(t *testing.T) {
		rToken, aToken, err := service.RegisterUser(
			context.TODO(),
			21,
			"juan diego bermudez",
			"male",
			"d@gmail.com",
			"actuacion.1",
		)
		if err != nil{
			t.Errorf("Expected no error and got %s", err.Error())
		}
		if rToken == ""{
			t.Error("Expected refresh token but not gotten")
		}
		if aToken == ""{
			t.Error("Expected access token but not gotten")
		}
	})

	t.Run("should receive already exists", func(t *testing.T) {
		_, _, err := service.RegisterUser(
			context.TODO(),
			21,
			"juan diego bermudez",
			"male",
			"d1@gmail.com",
			"actuacion.1",
		)
		if err != domain.UserWithEmailAlreadyExists{
			t.Error("Expected user already exists but got nothing")
		}
	})
}

type UserRepoMock struct{}
func (u * UserRepoMock) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error){
	if email == "d1@gmail.com"{
		return &domain.UserAuth{}, nil
	}
	return nil, domain.UserWithEmailAlreadyExists
}
func (u * UserRepoMock) RegisterUser(ctx context.Context, auth domain.UserAuth, user domain.User) (uuid.UUID, error) {
	return uuid.New(), nil
}


type SessionRepoMock struct{}

func (s SessionRepoMock) CreateSession(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error {
	return nil
}