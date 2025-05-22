package entities

import (
	"time"
)

type User struct {
	Id          int32     `json:"id"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	Name        string    `json:"name" binding:"required"`
	Surname     string    `json:"surname" binding:"required"`
	Patronymic  string    `json:"patronymic"`
	Age         int32     `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
}

type UpdateUser struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int32  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
