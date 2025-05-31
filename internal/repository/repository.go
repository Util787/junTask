package repository

import (
	"github.com/Util787/junTask/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetAllUsers(limit, offset int, name, surname, patronymic, gender string) ([]entities.User, error)
	CreateUser(params entities.User) (entities.User, error)
	ExistByFullName(params entities.FullName) (bool, error)
	ExistById(id int32) (bool, error)
	GetUserById(id int32) (entities.User, error)
	UpdateUser(id int32,params entities.UpdateUserParams) error
	DeleteUser(id int32) error
}

type Repository struct {
	UserRepository UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
