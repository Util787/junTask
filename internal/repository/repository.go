package repository

import (
	"context"

	"github.com/Util787/junTask/entities"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetAllUsers(pageSize, page int, name, surname, patronymic, gender string) (users []entities.User, totalCount int,err error)
	CreateUser(params entities.User) (entities.User, error)
	ExistByFullName(params entities.FullName) (bool, error)
	ExistById(id int32) (bool, error)
	GetUserById(id int32) (entities.User, error)
	UpdateUser(id int32, params entities.UpdateUserParams) error
	DeleteUser(id int32) error
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value any) error

	// Use reference for dest, otherwise you'll get an empty struct
	//
	// Example: Get(context.Background(),"key",&user)
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

type Repository struct {
	UserRepository  UserRepository
	RedisRepository RedisRepository
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		UserRepository:  NewUserRepository(db),
		RedisRepository: NewRedisRepository(redis),
	}
}
