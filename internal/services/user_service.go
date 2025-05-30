package service

import (
	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/repository"
)

type UserService struct {
	userRepo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{userRepo: repo}
}

func (u *UserService) CreateUser(params entities.User) (entities.User, error) {
	return u.userRepo.CreateUser(params)
}

func (u *UserService) GetAllUsers(limit, offset int, name, surname, patronymic, gender string) ([]entities.User, error) {
	return u.userRepo.GetAllUsers(limit, offset, name, surname, patronymic, gender)
}

func (u *UserService) ExistByFullName(params entities.FullName) (bool, error) {
	// its wrong to check existance if patronymic is null because people might have same names and surnames
	if params.Patronymic == "" {
		return false, nil
	}
	return u.userRepo.ExistByFullName(params)
}

func (u *UserService) ExistById(id int32) (bool, error) {
	return u.userRepo.ExistById(id)
}

func (u *UserService) GetUserById(id int32) (entities.User, error) {
	return u.userRepo.GetUserById(id)
}

func (u *UserService) UpdateUser(id int32, params entities.UpdateUserParams) error {
	return u.userRepo.UpdateUser(id, params)
}

func (u *UserService) DeleteUser(id int32) error {
	return u.userRepo.DeleteUser(id)
}
