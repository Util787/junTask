package entities

type User struct {
	Name        string `json:"name" binding:"required"`
	Surname     string `json:"surname" binding:"required"`
	Patronymic  string `json:"patronymic"`
	Age         int32  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type UpdateUser struct {
	Name        string `json:"name" binding:"omitempty"`
	Surname     string `json:"surname" binding:"omitempty"`
	Patronymic  string `json:"patronymic" binding:"omitempty"`
	Age         int32  `json:"age" binding:"omitempty"`
	Gender      string `json:"gender" binding:"omitempty"`
	Nationality string `json:"nationality" binding:"omitempty"`
}
