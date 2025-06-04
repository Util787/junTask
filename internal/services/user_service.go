package service

import (
	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/repository"
)

// TODO: implement validation logic from handlers here
type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (u *userService) CreateUser(params entities.User) (entities.User, error) {
	return u.userRepo.CreateUser(params)
}

func (u *userService) GetAllUsers(pageSize, page int, name, surname, patronymic, gender string) ([]entities.User, int, error) {
	return u.userRepo.GetAllUsers(pageSize, page, name, surname, patronymic, gender)
}

func (u *userService) ExistByFullName(params entities.FullName) (bool, error) {
	// its wrong to check existance if patronymic is null because people might have same names and surnames
	if params.Patronymic == "" {
		return false, nil
	}
	return u.userRepo.ExistByFullName(params)
}

func (u *userService) ExistById(id int32) (bool, error) {
	return u.userRepo.ExistById(id)
}

func (u *userService) GetUserById(id int32) (entities.User, error) {
	return u.userRepo.GetUserById(id)
}

func (u *userService) UpdateUser(id int32, params entities.UpdateUserParams) error {
	return u.userRepo.UpdateUser(id, params)
}

func (u *userService) DeleteUser(id int32) error {
	return u.userRepo.DeleteUser(id)
}
