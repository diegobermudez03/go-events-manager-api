package app

import (
	"context"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type RolesService struct {
	rolesRepo domain.RolesRepo
}

func NewRolesService(rolesRepo domain.RolesRepo) domain.RolesSvc{
	return &RolesService{
		rolesRepo: rolesRepo,
	}
}

func (s *RolesService) GetRoleById(ctx context.Context, roleId uuid.UUID) (*domain.Role, error){
	return nil, nil
}