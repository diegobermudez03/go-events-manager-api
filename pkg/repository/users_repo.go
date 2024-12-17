package repository

import (
	"context"
	"database/sql"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type UsersPostgres struct {
	db *sql.DB
}

func NewUsersPostgres(db *sql.DB) *UsersPostgres{
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, email, created_at
		FROM users_auth
		WHERE email = $1
		LIMIT 1`,
		email,
	)
	userAuth := new(domain.UserAuth)
	if err := row.Scan(&userAuth.Id, &userAuth.Email, &userAuth.CreatedAt); err != nil{
		if err == sql.ErrNoRows{
			return nil, domain.UserDoesntExistError 
		}
		return nil, err 
	}
	return userAuth, nil
}

func (r *UsersPostgres) RegisterUser(ctx context.Context, auth domain.UserAuth, user domain.User) error {
	return nil
}