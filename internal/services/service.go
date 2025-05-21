package service

import (
	"context"

	"github.com/Util787/junTask/internal/database"
	"github.com/Util787/junTask/internal/repository"
)

type User interface {
	Create(ctx context.Context, params database.CreateUserParams) (database.User, error)
	GetAll(ctx context.Context) ([]database.User, error)
}

type Service struct {
	UserService User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository),
	}
}
