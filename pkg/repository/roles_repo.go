package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type RolesPostgres struct {
	db	*sql.DB
}


func NewRolesPostgres(db *sql.DB) domain.RolesRepo{
	return &RolesPostgres{
		db: db,
	}
}

func (r *RolesPostgres) CreateRoleIfNotExists(ctx context.Context, role domain.Role) error{
	//check if the role already exists
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id FROM roles WHERE name = $1`,
		role.Name,
	)
	var id uuid.UUID
	if err := row.Scan(&id); err == nil{
		return nil
	}else if !errors.Is(err, sql.ErrNoRows){
		return err
	}

	//create the new role
	result, err := r.db.ExecContext(
		ctx, 
		`INSERT INTO roles(id, name)
		VALUES($1, $2)`,
		role.Id, role.Name,
	)
	if err != nil{
		return err
	}
	if num, err := result.RowsAffected(); num == 0 || err != nil{
		return errors.New("")
	}

	//create the permissions
	for _, permission := range role.Permissions{
		permissionId, err := r.CreatePermissionIfNotExists(ctx, permission)
		if err != nil{
			return err
		}
		intId := uuid.New()
		result, err := r.db.ExecContext(
			ctx,
			`INSERT INTO role_permissions(id, rolesid, permissionsid)
			VALUES($1, $2, $3)`,
			intId, role.Id, permissionId,
		)
		if err != nil{
			return err
 		}
		if num, err := result.RowsAffected(); num == 0 || err != nil{
			return errors.New("")
		}
	}
	return nil
}

func (r *RolesPostgres) CreatePermissionIfNotExists(ctx context.Context, permission string) (uuid.UUID, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id
		FROM permissions WHERE name = $1`,
		permission,
	)
	var id uuid.UUID
	if err :=row.Scan(&id); errors.Is(err, sql.ErrNoRows){
		id = uuid.New()
		result, err := r.db.ExecContext(
			ctx,
			`INSERT INTO permissions(id, name)
			VALUES($1, $2)`,
			id, permission,
		)
		if err != nil{
			return id, err
		}
		if num, err := result.RowsAffected(); num == 0 || err != nil{
			return id, errors.New("")
		}
		return id, nil

	}else if err != nil{
		return id, err
	}
	return id, nil
}

func (r *RolesPostgres) GetRoleByName(ctx context.Context, roleName string) (*domain.DataModelRole, error){
	return nil, nil
}

func (r *RolesPostgres) GetRoleIdByName(ctx context.Context, roleName string) (uuid.UUID, error){
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id FROM roles WHERE name = $1`,
		roleName,
	)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil{
		return id, domain.ErrInternal
	}
	return id, nil
}

func (r *RolesPostgres) GetRoleById(ctx context.Context, roleId uuid.UUID) (*domain.DataModelRole, error){
	row := r.db.QueryRowContext(
		ctx, 
		`SELECT id, name
		FROM roles WHERE id = $1`,
		roleId,
	)
	dataModel := new(domain.DataModelRole)
	if err := row.Scan(&dataModel.Id, &dataModel.Name); err != nil{
		return nil, domain.ErrInternal
	}
	return dataModel, nil
}