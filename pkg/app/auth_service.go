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

type AuthService struct {
	authRepo				domain.AuthRepo
	usersRepo 				domain.UsersRepo
	sessionsRepo 			domain.SessionsRepo
	eventsRepo 				domain.EventsRepo
	rolesRepo 				domain.RolesRepo
	tokensLifeHours	 		int64
	accessTokensLife 		int64
	jwtSecret 				string
	genderMap				map[string]string
}

func NewAuthService(
	authRepo 			domain.AuthRepo,
	usersRepo 			domain.UsersRepo, 
	sessionsRepo 		domain.SessionsRepo, 
	eventsRepo 			domain.EventsRepo,
	rolesRepo 			domain.RolesRepo,
	tokensLife 			int64, 
	accessTokensLife 	int64,
	jwtSecret 			string,
) domain.AuthSvc {
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
		rolesRepo: rolesRepo,
		eventsRepo: eventsRepo,
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
	session, err := s.sessionsRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil{
		return "", err
	}
	if time.Now().After(session.ExpiresAt) {
		 _ = s.sessionsRepo.DeleteSessionById(ctx, session.Id)
		 return "", domain.ErrExpiredSession
	}
	user, err := s.authRepo.GetUserAuthById(ctx, session.UserId)
	if err != nil{
		return "", err
	}
	accessToken, err := s.generateAccessToken(*user)
	if err != nil{
		return "", err
	}
	return accessToken, nil
}


func (s *AuthService) CheckAuthEvent(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, neededPermissions []string) error{
	participation, err := s.eventsRepo.GetParticipation(ctx, eventId, userId)
	if err != nil{
		return err 
	}
	permissions, err := s.rolesRepo.GetRolePermissions(ctx, participation.RoleId)
	if err != nil{
		return err 
	}
	permissionsMap := map[string]string{}
	for _, perm := range permissions{
		permissionsMap[perm.Name] = perm.Name
	}
	//check if all permissions needed are in the role
	for _, needed := range neededPermissions{
		if _, ok := permissionsMap[needed]; !ok{

			return domain.ErrUnathorized
		}
	}
	return nil
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
		UserId: user.Id,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.tokensLifeHours * 3600)),
	}
	//	create session in database
	err = s.sessionsRepo.CreateSession(ctx, session)
	if err != nil{
		return "", err
	}
	return token, nil
}

//this function doesn't validate refresh token, that must be handled outside
func (s *AuthService) generateAccessToken(user domain.UserAuth) (string, error){
	customClaims := domain.CustomJWTClaims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessTokensLife*int64(time.Second)))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "", err 
	}
	return tokenString, nil
}