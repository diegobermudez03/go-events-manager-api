package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, full_name, birth_date, gender, created_at
		FROM users
		WHERE id = $1`,
		userId,
	)
	user := new(domain.User)
	if err := row.Scan(&user.Id, &user.FullName, &user.BirthDate, &user.Gender, &user.CreatedAt); err != nil{
		return nil, domain.ErrInternal
	}
	return user, nil
}

func (r *UsersPostgres) GetUsers(ctx context.Context, filters domain.UsersFilters) ([]domain.User, error){
	rows, err := r.db.QueryContext(
		ctx, 
		`SELECT id, full_name, birth_date, gender, created_at
		FROM users
		WHERE
		($1::TEXT IS NULL OR (full_name ILIKE ('%' || $1::TEXT || '%') OR gender ILIKE  ('%' || $1::TEXT || '%')))
		LIMIT $2::INTEGER
		OFFSET $3::INTEGER`,
		filters.Text, filters.Limit, filters.Offset,
	)
	if err != nil{
		log.Println(err.Error())
		return nil, domain.ErrInternal
	}
	users := []domain.User{}
	for rows.Next(){
		user := domain.User{}
		if err := rows.Scan(&user.Id, &user.FullName, &user.BirthDate, &user.Gender, &user.CreatedAt); err != nil{
			log.Println(err.Error())
			return nil, domain.ErrInternal
		}
		users = append(users, user)
	}
	return users, nil
}