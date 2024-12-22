package domain

import "github.com/google/uuid"

type DataModelParticipation struct {
	Id 			uuid.UUID
	UserId		uuid.UUID 
	EventId 	uuid.UUID
	RoleId 		uuid.UUID
}

type DataModelRole struct{
	Id 		uuid.UUID
	Name 	string 
}

type DataModelPermission struct{
	Id 		uuid.UUID
	Name 	string
}