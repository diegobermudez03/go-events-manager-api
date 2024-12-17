package app

import (
	"context"
	"strings"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/helpers"
)

type AuthService struct {
	usersRepo domain.UsersRepo
	genderMap	map[string]string
}

func NewAuthService(usersRepo domain.UsersRepo) *AuthService {
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
		genderMap: genderMap,
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
	if err := s.usersRepo.RegisterUser(ctx, userAuth, user); err != nil{
		return "", "", nil
	}

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
	
	return "", nil
}

//this function doesn't validate refresh token, that must be handled outside
func (s *AuthService) generateAccessToken(ctx context.Context, user domain.UserAuth) (string, error){
	return "", nil
}