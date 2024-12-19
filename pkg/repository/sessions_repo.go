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

func (r *SessionsPotsgres) CreateSession(ctx context.Context, session domain.Session) error {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO sessions(id, refresh_token, created_at, expires_at, user_id)
		VALUES($1, $2, $3, $4, $5)`,
		session.Id, session.Token, session.CreatedAt, session.ExpiresAt, session.UserId,
	)
	if err != nil{
		log.Println(err.Error())
		return domain.ErrInternal 
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return domain.ErrInternal
	}
	return nil
}

func (r *SessionsPotsgres) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, refresh_token, created_at, expires_at, user_id
		FROM sessions WHERE refresh_token = $1`,
		token,
	)
	session := new(domain.Session)
	err := row.Scan(&session.Id, &session.Token, &session.CreatedAt, &session.ExpiresAt, &session.UserId)
	if err != nil{
		return nil ,domain.ErrInternal
	}
	return session, nil
}

func (r *SessionsPotsgres) DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM sessions WHERE id = $1`,
		sessionId,
	)
	if err != nil{
		return domain.ErrInternal
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return domain.ErrInternal
	}
	return nil
}