package domain

import (
	"time"

	"github.com/google/uuid"
)

type InvitationsEventBus interface {
	Suscribe(eventId uuid.UUID) <-chan InvitationEvent
	Unsuscribe(eventId uuid.UUID, channel <-chan InvitationEvent)
	Publish(eventId uuid.UUID, event InvitationEvent)
}

type InvitationEvent struct{
	Email 		string 		`json:"email"`
	Time 		time.Time 	`json:"time"`
}
