package service

import (
	"context"

	"github.com/Util787/junTask/internal/database"
	"github.com/Util787/junTask/internal/repository"
)

type User interface {
	GetAll(ctx context.Context) ([]database.User, error)
	Create(ctx context.Context, params database.CreateUserParams) (database.User, error)
	ExistByFullName(ctx context.Context, params database.UserExistByFullNameParams) (bool, error)
	ExistById(ctx context.Context, id int32) (bool, error)
	GetUserById(ctx context.Context, id int32) (database.User, error)
	UpdateUser(ctx context.Context, params database.UpdateUserParams) error
	DeleteUser(ctx context.Context, id int32) error
}

type Service struct {
	UserService User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository),
	}
}
