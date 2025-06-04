package service

import (
	"context"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/repository"
)

type UserService interface {
	GetAllUsers(pageSize, page int, name, surname, patronymic, gender string) (users []entities.User, totalCount int,err error)
	CreateUser(params entities.User) (entities.User, error)
	ExistByFullName(params entities.FullName) (bool, error)
	ExistById(id int32) (bool, error)
	GetUserById(id int32) (entities.User, error)
	UpdateUser(id int32, params entities.UpdateUserParams) error
	DeleteUser(id int32) error
}

type RedisService interface {
	Set(ctx context.Context, key string, value any) error
	
	// Use reference for dest, otherwise you'll get an empty struct
	//
	// Example: Get(context.Background(),"key",&user)
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

type InfoRequestService interface {
	RequestAdditionalInfo(name string) (age int, gender string, nationality string, err error)
}

type Service struct {
	UserService        UserService
	RedisService       RedisService
	InfoRequestService InfoRequestService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService:        NewUserService(repos.UserRepository),
		RedisService:       NewRedisService(repos.RedisRepository),
		InfoRequestService: NewInfoRequestService(),
	}
}
