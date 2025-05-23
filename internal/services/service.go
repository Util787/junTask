package service

import (
	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/repository"
)

type User interface {
	GetAllUsers(limit, offset int, name, surname, patronymic, gender string) ([]entities.User, error)
	CreateUser(params entities.User) (entities.User, error)
	ExistByFullName(params entities.FullName) (bool, error)
	ExistById(id int32) (bool, error)
	GetUserById(id int32) (entities.User, error)
	UpdateUser(params entities.UpdateUserParams) error
	DeleteUser(id int32) error
}

type Service struct {
	UserService User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository),
	}
}
