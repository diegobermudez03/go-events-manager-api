package domain

import (
	"time"

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
