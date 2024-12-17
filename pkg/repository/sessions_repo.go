package repository

import (
	"context"
	"database/sql"
	"time"

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

func (r *SessionsPotsgres) CreateSession(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error {
	return nil
}