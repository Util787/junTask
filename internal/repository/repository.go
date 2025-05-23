package repository

import (
	"github.com/Util787/junTask/entities"
	"github.com/jmoiron/sqlx"
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

type Repository struct {
	UserRepository User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
