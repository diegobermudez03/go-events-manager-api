package app

import (
	"context"

	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
)

type UsersService struct {
	usersRepo domain.UsersRepo
}

func NewUsersService(usersRepo domain.UsersRepo) domain.UserSvc{
	return &UsersService{
		usersRepo: usersRepo,
	}
}

func (s *UsersService) GetUsers(ctx context.Context, filters ...domain.UsersFilter) ([]domain.User, error){
	usersFilters := domain.UsersFilters{}
	// apply filters and check limit and offset
	for _, f := range filters{
		f(&usersFilters)
	}

	if usersFilters.Limit == nil{
		usersFilters.Limit = new(int)
		*usersFilters.Limit = 100
	}
	if usersFilters.Offset == nil{
		usersFilters.Offset = new(int)
		*usersFilters.Offset = 0
	}

	//get users with filters
	return s.usersRepo.GetUsers(ctx, usersFilters)
}