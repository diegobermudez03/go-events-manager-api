package app

import (
	"context"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type EventsService struct{}

func NewEventsService() domain.EventsSvc{
	return &EventsService{}
}

func (s *EventsService) CreateEvent(ctx context.Context, event domain.CreateEventRequest, profilePic *[]byte, creatorId uuid.UUID) error {
	return nil
}