package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/diegobermudez03/go-events-manager-api/internal/utils"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
)


var (
	ErrNoToken = errors.New("no access token")
	ErrMalformedToken = errors.New("malformed token")
	ErrExpiredToken = errors.New("expired token")
)

type Middlewares struct {
	jwtSecret 	string
}

// context values
type UserId string
const UserIdKey UserId = "userId"

func NewMiddlewares(jwtSecret string) *Middlewares {
	return &Middlewares{
		jwtSecret: jwtSecret,
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
