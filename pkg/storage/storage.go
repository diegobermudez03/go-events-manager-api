package storage

import (
	"database/sql"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/repository"
)

type Storage struct {
	UsersRepo 		domain.UsersRepo
	AuthRepo 		domain.AuthRepo
	SessionsRepo 	domain.SessionsRepo
	RolesRepo		domain.RolesRepo
	EventsRepo 		domain.EventsRepo
	FilesRepo 		domain.FilesRepo
}

func NewPostgreStorage(db *sql.DB, filesRepo domain.FilesRepo) *Storage{
	return &Storage{
		UsersRepo: repository.NewUsersPostgres(db),
		AuthRepo: repository.NewAuthPostgres(db),
		SessionsRepo: repository.NewSessionsPostgres(db),
		RolesRepo : repository.NewRolesPostgres(db),
		EventsRepo: repository.NewEventsPostgres(db),
		FilesRepo: filesRepo,
	}
}