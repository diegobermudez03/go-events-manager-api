package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type SessionsPotsgres struct {
	db *sql.DB
}

func NewSessionsPostgres(db *sql.DB) *SessionsPotsgres{
	return &SessionsPotsgres{
		db: db,
	}
}

func (r *SessionsPotsgres) CreateSession(ctx context.Context, session domain.Session, userId uuid.UUID) error {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions(id, refresh_token, created_at, expires_at, user_id)
		VALUES($1, $2, $3, $4, $5)`,
		session.Id, session.Token, session.Created_at, session.Expires_at, userId,
	)
	if err != nil{
		log.Println(err.Error())
		return err 
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return err
	}
	return nil
}