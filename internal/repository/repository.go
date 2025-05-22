package repository

import (
	"context"
	"database/sql"

	"github.com/Util787/junTask/internal/database"
)

type User interface {
	GetAll(ctx context.Context) ([]database.User, error)
	Create(ctx context.Context, params database.CreateUserParams) (database.User, error)
	Exist(ctx context.Context, params database.UserExistsParams) (bool, error)
	GetUserById(ctx context.Context, id int32) (database.User, error)
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
