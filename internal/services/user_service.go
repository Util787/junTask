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
	createdUser, err := u.userRepo.Create(ctx, params)
	if err != nil {
		return database.User{}, err
	}
	return createdUser, nil
}

func (u *UserService) GetAll(ctx context.Context) ([]database.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *UserService) Exist(ctx context.Context, params database.UserExistsParams) (bool, error) {
	// its wrong to check if patronymic is null because people might have same names and surnames
	if params.Patronymic.String == "" {
		return false, nil
	}
	return u.userRepo.Exist(ctx, params)
}

func (u *UserService) GetUserById(ctx context.Context, id int32) (database.User, error) {
	return u.userRepo.GetUserById(ctx, id)
}
