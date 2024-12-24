package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

// Using joins due to clarity and readability, however using EXISTS is better for performance, in this case, I'll stay with join
// In this case since I'm using Postgres, would be recommended to create Indexes for FK, maybe I'll add that to migrations later
func (r *EventsPostgres) GetParticipations(ctx context.Context, filters domain.ParticipationFilters) ([]domain.DataModelParticipation, error){
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT p.id, p.userid, p.eventid, p.roleid
		FROM participants p
		INNER JOIN roles r ON r.id = p.roleid 
		INNER JOIN events e ON e.id = p.eventid
		INNER JOIN users u ON u.id = p.userid 
		WHERE 
		($1::UUID IS NULL OR p.userid = $1::UUID) AND
		($2::TEXT IS NULL OR r.name = $2::TEXT) AND
		($3::UUID IS NULL OR p.eventid = $3::UUID)
		LIMIT $4::INTEGER
		OFFSET COALESCE($5::INTEGER,0)`,
		filters.UserId, filters.RoleName, filters.EventId, filters.Limit, filters.Offset,
	)
	if err != nil{
		return nil, domain.ErrInternal
	}
	participations := []domain.DataModelParticipation{} 

	for rows.Next(){
		part := domain.DataModelParticipation{}
		if err :=rows.Scan(&part.Id, &part.UserId, &part.EventId, &part.RoleId); err != nil{
			return nil, domain.ErrInternal
		}
		participations = append(participations, part)
	}
	return participations, nil
}

func (r *EventsPostgres) GetEventById(ctx context.Context, eventId uuid.UUID) (*domain.Event, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, description, starts_at, ends_at, profile_pic_url, address, created_at
		FROM events
		WHERE id = $1`,
		eventId,
	)
	event := new(domain.Event)
	if err := row.Scan(&event.Id, &event.Name, &event.Description, &event.StartsAt, &event.EndsAt, &event.ProfilePicUrl, &event.Address, &event.CreatedAt); err != nil{
		return nil, domain.ErrInternal
	}
	return event, nil
}

func (r *EventsPostgres) GetParticipation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID)(*domain.DataModelParticipation, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, userid, eventid, roleid
		FROM participants
		WHERE userid = $1 AND eventid = $2`,
		userId, eventId,
	)

	dataModel := new(domain.DataModelParticipation)

	if err := row.Scan(&dataModel.Id, &dataModel.UserId, &dataModel.EventId, &dataModel.RoleId); err != nil{
		return nil, domain.ErrNoParticipationFound
	}
	return dataModel, nil
}

func (r *EventsPostgres) CreateInvitation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error{
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO invitations(id, userid, eventid, created_at)
		VALUES($1, $2, $3, $4)`,
		uuid.New(), userId, eventId, time.Now(),
	)
	if err != nil{
		return domain.ErrInternal
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return domain.ErrInternal
	}
	return nil
}

func (r *EventsPostgres) CheckInvitation(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) (bool, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id
		FROM invitations
		WHERE userid = $1 AND eventId = $2`,
		userId, eventId,
	)

	var id uuid.UUID
	if err := row.Scan(&id); errors.Is(err, sql.ErrNoRows){
		return false, nil 
	}else if err != nil{
		return false, err 
	}
	return true, nil
}