package mock

import (
	"context"
	"errors"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

/*
	RIGHT NOW I'M CREATING THE MOCK DATA HERE IN THE INIT FUNCTION, HOWEVER, THE IDEAL THING IS
	TO DO IT OUTSIDE, MAYBE IN A JSON FILE, AND THEN HERE SIMPLY READ THE DATA, THIS WAY,
	WE CAN HAVE MULTIPLE JSON FILES, EACH ONE CORRESPONDING TO A DIFFERENT SCENARIO, I'LL LEAVE THIS
	TO A FUTURE CHANGE
*/


var partDataModels []domain.DataModelParticipation
var users 	[]domain.User
var events 	[]domain.Event
var roles	[]domain.Role

func init(){
	//users ID's
	userId1 := uuid.New(); userId2 := uuid.New(); userId3 := uuid.New(); /*userId4 := uuid.New();
	userId5 := uuid.New(); userId6 := uuid.New(); userId7 := uuid.New(); userId8 := uuid.New();
	userId9 := uuid.New(); userId10 := uuid.New(); userId11 := uuid.New(); userId12 := uuid.New();*/

	//events ID's
	eventId1 := uuid.New(); eventId2 := uuid.New(); eventId3 := uuid.New(); eventId4 := uuid.New();
	eventId5 := uuid.New(); /*eventId6 := uuid.New(); eventId7 := uuid.New(); eventId8 := uuid.New();
	eventId9 := uuid.New(); eventId10 := uuid.New(); eventId11 := uuid.New(); eventId12 := uuid.New();*/

	//roles creation
	roleId1 := uuid.New(); roleId2 := uuid.New(); roleId3 := uuid.New();
	roles = make([]domain.Role, 3)
	roles[0] = domain.Role{roleId1, "Creator", []string{}}
	roles[1] = domain.Role{roleId2, "Administrator", []string{}}
	roles[2] = domain.Role{roleId3, "Participant", []string{}}

	//events creation
	events = make([]domain.Event, 5)
	events[0] = domain.Event{eventId1, "Comic con", "for comics fans", time.Now(), time.Now(), "///", "san diego", time.Now()}
	events[1] = domain.Event{eventId2, "Corferias", "for corferias fans", time.Now(), time.Now(), "///", "Colombia", time.Now()}
	events[2] = domain.Event{eventId3, "Expojaveriana", "for university", time.Now(), time.Now(), "///", "Bogota", time.Now()}
	events[3] = domain.Event{eventId4, "Feria del carro", "for cars", time.Now(), time.Now(), "///", "Bogota", time.Now()}
	events[4] = domain.Event{eventId5, "Filbo", "for books", time.Now(), time.Now(), "///", "Medellin", time.Now()}

	//users creation
	users = make([]domain.User, 3)
	users[0] = domain.User{userId1, "Juan DIego", time.Now(), "male", time.Now()}
	users[1] = domain.User{userId2, "Carlos", time.Now(), "male", time.Now()}
	users[2] = domain.User{userId3, "Gabriela", time.Now(), "female", time.Now()}

	//participations data models creation
	partDataModels = make([]domain.DataModelParticipation, 6)
	partDataModels[0] = domain.DataModelParticipation{uuid.New(), userId1, eventId1, roleId1}
	partDataModels[1] = domain.DataModelParticipation{uuid.New(), userId1, eventId2, roleId1}
	partDataModels[2] = domain.DataModelParticipation{uuid.New(), userId1, eventId3, roleId3}
	partDataModels[3] = domain.DataModelParticipation{uuid.New(), userId2, eventId1, roleId2}
	partDataModels[4] = domain.DataModelParticipation{uuid.New(), userId2, eventId5, roleId2}
	partDataModels[5] = domain.DataModelParticipation{uuid.New(), userId3, eventId1, roleId3}
}


////////////////	AUTH REPO
type AuthRepoMock struct{}
func (u *AuthRepoMock) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error) {
	if email == "d1@gmail.com" {
		return &domain.UserAuth{}, nil
	}
	return nil, domain.ErrUserWithEmailAlreadyExists
}

func (u *AuthRepoMock) GetUserAuthById(ctx context.Context, id uuid.UUID) (*domain.UserAuth, error) {
	return nil, domain.ErrUserWithEmailAlreadyExists
}

func (u *AuthRepoMock) RegisterUser(ctx context.Context, auth domain.UserAuth) error {
	return nil
}


////////////////	USERS REPO
type UsersRepoMock struct{}
func (u *UsersRepoMock) CreateUser(ctx context.Context, user domain.User) error {
	return nil
}
func (u *UsersRepoMock) GetUserById(ctx context.Context, userId uuid.UUID) (*domain.User, error){
	for _, user := range users{
		if userId == user.Id{
			return &user, nil  
		}
	}
	return nil, errors.New("")
}


////////////////	SESSIONS REPO
type SessionRepoMock struct{}

func (s SessionRepoMock) CreateSession(ctx context.Context, session domain.Session) error {
	return nil
}
func (r *SessionRepoMock) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error){
	return nil, nil
}

func (r *SessionRepoMock) DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error{
	return nil
}


////////////////	ROLES REPO
type RolesRepoMock struct{}

func (r *RolesRepoMock) CreateRoleIfNotExists(ctx context.Context, role domain.Role) error{
	return nil
}

func (r *RolesRepoMock) CreatePermissionIfNotExists(ctx context.Context, permission string) (uuid.UUID, error){
	var id uuid.UUID
	return id, nil
}

func (r *RolesRepoMock) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error){
	return nil, nil
}

func (r *RolesRepoMock) GetRoleIdByName(ctx context.Context, roleName string) (uuid.UUID, error){
	var id uuid.UUID
	return id, nil
}

func (r *RolesRepoMock) GetRoleById(ctx context.Context, roleId uuid.UUID) (*domain.Role, error){
	for _, role := range roles{
		if roleId == role.Id{
			return &role, nil  
		}
	}
	return nil, errors.New("")
}


////////////////	EVENTS REPO
type EventsRepoMock struct{}

func (r *EventsRepoMock) CreateEvent(ctx context.Context, event domain.Event) error {
	return nil
}

func (r *EventsRepoMock) CreateParticipant(ctx context.Context, userId uuid.UUID, eventId uuid.UUID, roleId uuid.UUID) error{
	return nil
}

func (r *EventsRepoMock) GetParticipations(ctx context.Context, filters domain.ParticipationFilters) ([]domain.DataModelParticipation, error){
	return partDataModels, nil
}

func (r *EventsRepoMock) GetEventById(ctx context.Context, eventId uuid.UUID) (*domain.Event, error){
	for _, event := range events{
		if eventId == event.Id{
			return &event, nil  
		}
	}
	return nil, errors.New("")
}



////////////////	FILES REPO


type FilesRepoMock struct {}

func (r *FilesRepoMock) StoreImage(ctx context.Context, image *[]byte, group string, imageName string, imageType string, path string ) (string, error){
	return "", nil
}

func (r *FilesRepoMock) DeleteFile(ctx context.Context, group string,url string) error{
	return nil
}

func (r *FilesRepoMock) getPath(group string, path string)(string, error){
	return "", nil
}