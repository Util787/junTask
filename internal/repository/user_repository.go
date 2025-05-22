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

func (u *UserRepository) GetAllUsers(ctx context.Context, params database.GetAllUsersParams) ([]database.User, error) {
	return u.dbQueries.GetAllUsers(ctx, params)
}

func (u *UserRepository) CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	return u.dbQueries.CreateUser(ctx, params)
}

func (u *UserRepository) ExistByFullName(ctx context.Context, params database.UserExistByFullNameParams) (bool, error) {
	return u.dbQueries.UserExistByFullName(ctx, params)
}

func (u *UserRepository) ExistById(ctx context.Context, id int32) (bool, error) {
	return u.dbQueries.UserExistById(ctx, id)
}

func (u *UserRepository) GetUserById(ctx context.Context, id int32) (database.User, error) {
	return u.dbQueries.GetUserById(ctx, id)
}

func (u *UserRepository) UpdateUser(ctx context.Context, params database.UpdateUserParams) error {
	return u.dbQueries.UpdateUser(ctx, params)
}

func (u *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	return u.dbQueries.DeleteUser(ctx, id)
}
