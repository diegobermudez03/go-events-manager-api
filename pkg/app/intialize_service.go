package app

import (
	"context"
	"slices"
	"time"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/google/uuid"
)

type InitializeService struct {
	rolesRepo 	domain.RolesRepo
}

func NewInitializeService(rolesRepo domain.RolesRepo) domain.InitializeSvc {
	return &InitializeService{
		rolesRepo: rolesRepo,
	}
}

func (s *InitializeService) RegisterRoles() error {
	for k, val := range domain.RolesPermissions{
		role := domain.Role{
			Id: uuid.New(),
			Name: k,
			Permissions: slices.Clone(val),
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.rolesRepo.CreateRoleIfNotExists(ctx, role)
		if err != nil{
			return err
		}
	}
	return nil
}