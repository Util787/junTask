package repository

import (
	"context"

	"github.com/Util787/junTask/internal/database"
)

type UserRepository struct {
	dbQueries *database.Queries
}

func NewUserRepository(dbQueries *database.Queries) *UserRepository {
	return &UserRepository{dbQueries: dbQueries}
}

func (u *UserRepository) GetAll(ctx context.Context) ([]database.User, error) {
	allUsers, err := u.dbQueries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (u *UserRepository) Create(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	allUsers, err := u.dbQueries.CreateUser(ctx, params)
	if err != nil {
		return database.User{}, err
	}
	return allUsers, nil
}
