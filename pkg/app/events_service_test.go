package app

import (
	"context"
	"testing"

	"github.com/diegobermudez03/go-events-manager-api/pkg/mock"
	"github.com/google/uuid"
)

func TestEventsService(t *testing.T) {
	usersMock := mock.UsersRepoMock{}
	eventsMock := mock.EventsRepoMock{}
	rolesMock := mock.RolesRepoMock{}
	filesMock := mock.FilesRepoMock{}
	authMock := mock.AuthRepoMock{}
	emailMock := mock.EmailMock{}
	invitationsBus := NewInvitationsEventBus()

	eventsService := NewEventsService(
		&eventsMock,
		&rolesMock,
		&usersMock,
		&filesMock,
		&authMock,
		&emailMock,
		invitationsBus,
	)
	t.Run("Should sucesfully retrieve the participations", func(t *testing.T) {
		participations, err := eventsService.GetParticipationsOfUser(context.Background(), uuid.New())
		if err != nil{
			t.Error("Error on Event service Test")
		}
		t.Logf("Retrieved: %v", participations)
	})
}