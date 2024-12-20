package domain

import (
	"context"
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

type AuthSvc interface {
	RegisterUser(ctx context.Context,age int, fullName, gender, email, password string) (string, string, error)
	LoginUser(ctx context.Context, email string, password string) (string, string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type UserSvc interface {
}

type EventsSvc interface{
	CreateEvent(ctx context.Context, event CreateEventRequest, profilePic *[]byte, creatorId uuid.UUID) error 
}
//	EventsSvc requests
type CreateEventRequest struct{
	Name 		string 
	Description string 
	StartsAt 	time.Time
	EndsAt		time.Time
	Address 	string 
}

type InitializeSvc interface {
	RegisterRoles() error 
}