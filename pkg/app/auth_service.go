package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const userIdKey string = "userId"

type AuthService struct {
	authRepo				domain.AuthRepo
	usersRepo 				domain.UsersRepo
	sessionsRepo 			domain.SessionsRepo
	tokensLifeHours	 		int64
	accessTokensLife 		int64
	jwtSecret 				string
	genderMap				map[string]string
}

func NewAuthService(
	authRepo 		domain.AuthRepo,
	usersRepo 		domain.UsersRepo, 
	sessionsRepo 	domain.SessionsRepo, 
	tokensLife 		int64, 
	accessTokensLife int64,
	jwtSecret 		string,
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
		authRepo: authRepo,
		usersRepo: usersRepo,
		sessionsRepo: sessionsRepo,
		tokensLifeHours: tokensLife,
		genderMap: genderMap,
		accessTokensLife: accessTokensLife,
		jwtSecret: jwtSecret,
	}
}

//////////////////////////////////////////////////////////////////////////////////////
////////////					PUBLIC METHODS				/////////////////////////
////////////////////////////////////////////////////////////////////////////////////

func (s *AuthService) RegisterUser(ctx context.Context,age int,fullName, gender,email, password string) (string, string, error) {
	//check if user already exists
	_, err := s.authRepo.GetUserAuthByEmail(ctx, email)
	if err == nil{
		return "", "", domain.ErrUserWithEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserDoesntExist){
		return "", "", domain.ErrInternal 
	}

	//hash password and convert to correct parameters
	hash, _ := helpers.HashPassword(password)
	birthDate := time.Now().AddDate(-age, 0,0)
	gender, found := s.genderMap[strings.ToLower(gender)]; 
	if !found{
		return "", "", domain.ErrInvalidParametersError
	}

	//generate new ID and create user in auth and user
	id := uuid.New()
	user := domain.User{
		Id: id,
		FullName: fullName,
		BirthDate: birthDate,
		Gender: gender,
	}
	userAuth := domain.UserAuth{
		Id: id,
		Email: email,
		Hash: hash,
		CreatedAt: time.Now(),
	}
	if err := s.authRepo.RegisterUser(ctx, userAuth); err != nil{
		return "", "", err
	}
	if err := s.usersRepo.CreateUser(ctx, user); err != nil {
		return "", "", err
	}

	//creating refresh and access token
	refreshToken, err := s.generateRefreshToken(ctx, userAuth)
	if err != nil{
		return "", "", domain.ErrInternal 
	}
	accessToken, err := s.generateAccessToken(userAuth)
	if err != nil{
		return "", "", domain.ErrInternal 
	}
	return refreshToken, accessToken, nil
}


func (s *AuthService) LoginUser(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.authRepo.GetUserAuthByEmail(ctx, email)
	if err != nil{
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash) ,[]byte(password)); err != nil{
		return "", "", domain.ErrIncorrectPassword
	}
	
	//	if correct was correct, then we generate tokens
	refreshToken, err := s.generateRefreshToken(ctx, *user)
	if err != nil{
		return "", "", domain.ErrInternal
	}
	accessToken, err := s.generateRefreshToken(ctx, *user)
	if err != nil{
		return "", "", domain.ErrInternal
	}
	return refreshToken, accessToken, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}


///////////////////////////////////////////////////////////////////////////////////////
////////////					helpers, aux functions				//////////////////
/////////////////////////////////////////////////////////////////////////////////////

func (s *AuthService) generateRefreshToken(ctx context.Context, user domain.UserAuth) (string, error){
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil{
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)
	session := domain.Session{
		Id: uuid.New(),
		Token: token,
		Created_at: time.Now(),
		Expires_at: time.Now().Add(time.Second * time.Duration(s.tokensLifeHours * 3600)),
	}
	//	create session in database
	err = s.sessionsRepo.CreateSession(ctx, session, user.Id)
	if err != nil{
		return "", err
	}
	return token, nil
}

//this function doesn't validate refresh token, that must be handled outside
func (s *AuthService) generateAccessToken(user domain.UserAuth) (string, error){
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