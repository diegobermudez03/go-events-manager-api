package app

import (
	"context"
	"testing"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/mock"
)

func TestAuthService(t *testing.T) {

	authRepo := mock.AuthRepoMock{}
	usersRepo := mock.UsersRepoMock{}
	sessionmock := mock.SessionRepoMock{}

	service := NewAuthService(&authRepo,&usersRepo, &sessionmock, 1440, 600, "secret")

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
		if err != domain.ErrUserWithEmailAlreadyExists{
			t.Error("Expected user already exists but got nothing")
		}
	})
}
