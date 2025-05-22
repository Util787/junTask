package service

import (
	"context"

	"github.com/Util787/junTask/internal/database"
	"github.com/Util787/junTask/internal/repository"
)

type UserService struct {
	userRepo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{userRepo: repo}
}

func (u *UserService) Create(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	return u.userRepo.Create(ctx, params)
}

func (u *UserService) GetAll(ctx context.Context) ([]database.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *UserService) ExistByFullName(ctx context.Context, params database.UserExistByFullNameParams) (bool, error) {
	// its wrong to check if patronymic is null because people might have same names and surnames
	if params.Patronymic.String == "" {
		return false, nil
	}
	return u.userRepo.ExistByFullName(ctx, params)
}

func (u *UserService) ExistById(ctx context.Context, id int32) (bool, error) {
	return u.userRepo.ExistById(ctx, id)
}

func (u *UserService) GetUserById(ctx context.Context, id int32) (database.User, error) {
	return u.userRepo.GetUserById(ctx, id)
}

func (u *UserService) UpdateUser(ctx context.Context, params database.UpdateUserParams) error {
	return u.userRepo.UpdateUser(ctx, params)
}

func (u *UserService) DeleteUser(ctx context.Context, id int32) error {
	return u.userRepo.DeleteUser(ctx, id)
}
