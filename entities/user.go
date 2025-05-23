package entities

import (
	"time"
)

type User struct {
	Id          int32     `json:"id" db:"id"`
	Created_at  time.Time `json:"created_at" db:"created_at"`
	Updated_at  time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Surname     string    `json:"surname" db:"surname" binding:"required"`
	Patronymic  string    `json:"patronymic" db:"patronymic"`
	Age         int       `json:"age" db:"age"`
	Gender      string    `json:"gender" db:"gender"`
	Nationality string    `json:"nationality" db:"nationality"`
}

type FullName struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
}

type UpdateUserParams struct {
	Id          int32     `json:"id"`
	Updated_at  time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
}
