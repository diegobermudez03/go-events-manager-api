package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type CustomJWTClaims struct{
	UserId 		uuid.UUID	`json:"userId"`
	jwt.RegisteredClaims
}

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

type Session struct{
	Id 			uuid.UUID
	Token 		string
	UserId		uuid.UUID
	CreatedAt	time.Time 
	ExpiresAt 	time.Time
}

type Role struct{
	Id 			uuid.UUID
	Name 		string 
	Permissions []string
}

type Event struct{
	Id 				uuid.UUID
	Name 			string 
	Description 	string
	StartsAt		time.Time 
	EndsAt 			time.Time 
	ProfilePicUrl 	string 
	Address 		string 
	CreatedAt 		time.Time 
}

type Participation struct{
	Id 			uuid.UUID
	Event		*Event
	RoleName 	string 
	User 		*User 		
}

const EventsGroup = "events"