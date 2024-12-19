package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres{
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error){
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
			return nil, domain.ErrUserDoesntExist
		}
		log.Printf("Auth repo error %s", err.Error())
		return nil, domain.ErrInternal 
	}
	return userAuth, nil
}

func (r *AuthPostgres) GetUserAuthById(ctx context.Context, id uuid.UUID) (*domain.UserAuth, error){
	return nil, nil
}

func (r *AuthPostgres) RegisterUser(ctx context.Context, auth domain.UserAuth) error {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users_auth(id, email, hash, created_at)
		VALUES($1, $2, $3, $4)`,
		auth.Id, auth.Email, auth.Hash, auth.CreatedAt,
	)
	if err != nil{
		log.Printf("Auth repo error %s", err.Error())
		return domain.ErrInternal 
	}
	if number, err := result.RowsAffected(); number == 0 || err != nil{
		log.Printf("Auth repo error %s", err.Error())
		return domain.ErrInternal
	}
	return  nil
}