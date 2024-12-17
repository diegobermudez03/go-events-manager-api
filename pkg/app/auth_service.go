package app

import (
	"context"
	"crypto/rand"
	"strings"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/helpers"
	"github.com/golang-jwt/jwt/v5"
)

const userIdKey string = "userId"

type AuthService struct {
	usersRepo 		domain.UsersRepo
	sessionsRepo 	domain.SessionsRepo
	tokensLife 		int64
	accessTokensLife int64
	jwtSecret 		string
	genderMap		map[string]string
}

func NewAuthService(
	usersRepo domain.UsersRepo, 
	sessionsRepo domain.SessionsRepo, 
	tokensLife int64, 
	accessTokensLife int64,
	jwtSecret string,
) *AuthService {
	genderMap := map[string]string{
		"male" : "MALE",
		"female" : "FEMALE",
		"masculino" : "MALE",
		"femenino" : "FEMALE",
		"hombre" : "MALE",
		"mujer" : "FEMALE",
		"man" : "MALE",
		"woman" : "FEMALE",
	}
	return &AuthService{
		usersRepo: usersRepo,
		sessionsRepo: sessionsRepo,
		tokensLife: tokensLife,
		genderMap: genderMap,
		accessTokensLife: accessTokensLife,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context,age int,fullName, gender,email, password string) (string, string, error) {
	//check if user already exists
	_, err := s.usersRepo.GetUserAuthByEmail(ctx, email)
	if err == nil{
		return "", "", domain.UserWithEmailAlreadyExists
	}
	if err != domain.UserWithEmailAlreadyExists{
		return "", "", err 
	}

	//hash password and convert to correct parameters
	hash, _ := helpers.HashPassword(password)
	userAuth := domain.UserAuth{
		Email: email,
		Hash: hash,
		CreatedAt: time.Now(),
	}
	birthDate := time.Now().AddDate(-age, 0,0)
	gender, found := s.genderMap[strings.ToLower(gender)]; 
	if !found{
		return "", "", domain.InvalidParametersError
	}
	user := domain.User{
		FullName: fullName,
		BirthDate: birthDate,
		Gender: gender,
	}

	//register user
	id, err := s.usersRepo.RegisterUser(ctx, userAuth, user);
	if err != nil{
		return "", "", nil
	}
	userAuth.Id = id

	//creating refresh and access token
	refreshToken, err := s.generateRefreshToken(ctx, userAuth)
	if err != nil{
		return "", "", err 
	}
	accessToken, err := s.generateAccessToken(ctx, userAuth)
	if err != nil{
		return "", "", err
	}
	return refreshToken, accessToken, nil
}


func (s *AuthService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}

func (s *AuthService) generateRefreshToken(ctx context.Context, user domain.UserAuth) (string, error){
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil{
		return "", err
	}
	token := string(randomBytes)
	err = s.sessionsRepo.CreateSession(ctx, user.Id, token, time.Now().Add(time.Second * time.Duration(s.tokensLife)))
	if err != nil{
		return "", err
	}
	return token, nil
}

//this function doesn't validate refresh token, that must be handled outside
func (s *AuthService) generateAccessToken(ctx context.Context, user domain.UserAuth) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		userIdKey : user.Id,
		"exp" : time.Now().Add(time.Duration(s.accessTokensLife*int64(time.Second))).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "", err 
	}
	return tokenString, nil
}