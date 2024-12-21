package domain

import "github.com/google/uuid"

type DataModelParticipation struct {
	Id 			uuid.UUID
	UserId		uuid.UUID 
	EventId 	uuid.UUID
	RoleId 		uuid.UUID
}