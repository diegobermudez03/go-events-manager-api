package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) domain.AuthRepo{
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) GetUserAuthByEmail(ctx context.Context, email string) (*domain.UserAuth, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, email, hash, created_at
		FROM users_auth
		WHERE email = $1
		LIMIT 1`,
		email,
	)
	userAuth, err := r.rowToUserAuth(row)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, domain.ErrUserDoesntExist
		}
		return nil, domain.ErrInternal
	}
	return userAuth, nil
}

func (r *AuthPostgres) GetUserAuthById(ctx context.Context, id uuid.UUID) (*domain.UserAuth, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, email, hash, created_at
		FROM users_auth
		WHERE id = $1`,
		id,
	)
	userAuth, err := r.rowToUserAuth(row)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, domain.ErrUserDoesntExist
		}
		return nil, domain.ErrInternal
	}
	return userAuth, nil
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

func (r *AuthPostgres) rowToUserAuth(row *sql.Row) (*domain.UserAuth, error){
	user := new(domain.UserAuth)
	if err := row.Scan(&user.Id, &user.Email, &user.Hash, &user.CreatedAt); err != nil{
		return nil, err
	}
	return user, nil
}