package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


var (
	ErrNoToken = errors.New("no access token")
	ErrMalformedToken = errors.New("malformed token")
	ErrExpiredToken = errors.New("expired token")
)

type Middlewares struct {
	jwtSecret 	string
	authService domain.AuthSvc
}

// context values
type UserId string
const UserIdKey UserId = "userId"

func NewMiddlewares(jwtSecret string, authService domain.AuthSvc) *Middlewares {
	return &Middlewares{
		jwtSecret: jwtSecret,
		authService: authService,
	}
}

func (m *Middlewares) AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//get header
			tokenString := r.Header.Get("Authorization")
			if tokenString == ""{
				utils.WriteError(w, http.StatusBadRequest, ErrNoToken)
				return 
			}
			tokenString = strings.Split(tokenString, " ")[1]	//to remove the Bearer

			//parse jwt token
			token, err := jwt.ParseWithClaims(tokenString, &domain.CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(m.jwtSecret), nil
			})
			if err != nil{
				if errors.Is(err, jwt.ErrTokenExpired){
					utils.WriteError(w, http.StatusBadRequest, ErrExpiredToken)
					return
				}
				utils.WriteError(w, http.StatusBadRequest, ErrMalformedToken)
				return 
			}
			if !token.Valid{
				utils.WriteError(w, http.StatusBadRequest, ErrExpiredToken)
				return
			}
			claims := token.Claims.(*domain.CustomJWTClaims)
			id := claims.UserId
			ctx := context.WithValue(r.Context(), UserIdKey, id)
			r = r.WithContext(ctx)
			log.Printf("Succesfully validated %s", id)
			next.ServeHTTP(w, r)
		},
	)
}


func (m *Middlewares) EventAccessMiddleware(permissions ...string) func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(UserIdKey).(uuid.UUID)
			eventIdString := chi.URLParam(r, "eventId")
			if eventIdString == ""{
				utils.WriteError(w, http.StatusBadRequest, errors.New("missing event Id"))
				return 
			}
			eventId, err := uuid.Parse(eventIdString)
			if err != nil{
				utils.WriteError(w, http.StatusBadRequest, errors.New("invalid event Id"))
				return 
			}
			if err := m.authService.CheckAuthEvent(r.Context(), eventId, userId, permissions); err != nil{
				utils.WriteError(w, http.StatusUnauthorized, errors.New("user not authorized"))
				return
			}
			log.Println("succesfully validated event auth")
			ctx := context.WithValue(r.Context(), "eventId", eventId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}