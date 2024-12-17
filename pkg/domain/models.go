package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	FullName  string    `json:"fullName"`
	BirthDate time.Time `json:"birthDate"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAuth struct {
	Id        uuid.UUID
	Email     string
	Hash      string
	CreatedAt time.Time
}
