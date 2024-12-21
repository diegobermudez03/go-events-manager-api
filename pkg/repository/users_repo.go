package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type UsersPostgres struct {
	db *sql.DB
}

func NewUsersPostgres(db *sql.DB) domain.UsersRepo{
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) CreateUser(ctx context.Context, user domain.User) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users(id, full_name, birth_date, gender, created_at)
		VALUES($1, $2, $3, $4, $5)`,
		user.Id, user.FullName, user.BirthDate, user.Gender, user.CreatedAt,
	)
	if err != nil{
		return err 
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return errors.New("")
	}
	return nil
}

func (r *UsersPostgres) GetUserById(ctx context.Context, userId uuid.UUID) (*domain.User, error){
	return nil, nil
}