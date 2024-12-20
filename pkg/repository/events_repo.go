package repository

import (
	"context"
	"database/sql"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type EventsPostgres struct {
	db *sql.DB
}

func NewEventsPostgres(db *sql.DB) domain.EventsRepo{
	return &EventsPostgres{
		db: db,
	}
}

func (r *EventsPostgres) CreateEvent(ctx context.Context, event domain.Event) error{
	return nil
}