package repository

import (
	"context"
	"database/sql"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type EventsPostgres struct {
	db *sql.DB
}

func NewEventsPostgres(db *sql.DB) domain.EventsRepo{
	return &EventsPostgres{
		db: db,
	}
}

func (r *EventsPostgres) CreateEvent(ctx context.Context, event domain.Event) error {
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO events(id, name, description, starts_at, ends_at, profile_pic_url, address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.Id, event.Name, event.Description, event.StartsAt, event.EndsAt, event.ProfilePicUrl, 
		event.Address, event.CreatedAt,
	)
	if err != nil{
		return domain.ErrInternal
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return domain.ErrInternal
	}
	return nil
}

func (r *EventsPostgres) CreateParticipant(ctx context.Context, userId uuid.UUID, eventId uuid.UUID, roleId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO participants(id, userid, eventid, roleid)
		VALUES($1, $2, $3, $4)`,
		uuid.New(), userId, eventId, roleId,
	)
	if err != nil{
		return domain.ErrInternal
	}
	if num, err := result.RowsAffected(); err != nil || num == 0{
		return domain.ErrInternal
	}
	return nil
}