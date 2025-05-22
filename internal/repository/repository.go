package repository

import (
	"context"
	"database/sql"

	"github.com/Util787/junTask/internal/database"
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

type Repository struct {
	UserRepository User
}

func NewRepository(db *sql.DB) *Repository {
	dbQueries := database.New(db)
	return &Repository{
		UserRepository: NewUserRepository(dbQueries),
	}
}
