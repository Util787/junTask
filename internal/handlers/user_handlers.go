package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Util787/junTask/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	allUsers, err := h.services.UserService.GetAll(context.Background())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, allUsers)
}

func (h *Handler) createUser(c *gin.Context) {
	// using entities.User instead of database.User is important because of tags
	var user entities.User
	err := c.BindJSON(&user)
	if err != nil {
		return
	}

	age, gender, nationality := requestUserAdditionalInfo(c, user.Name)
	log.Println(user, age, gender, nationality)
	// params := database.CreateUserParams{
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// 	Name:        user.Name,
	// 	Surname:     user.Surname,
	// 	Age:         age,
	// 	Gender:      gender,
	// 	Nationality: nationality,
	// }
	// if user.Patronymic == "" {
	// 	params.Patronymic = sql.NullString{Valid: false}
	// } else {
	// 	params.Patronymic = sql.NullString{String: user.Patronymic, Valid: true}
	// }

	// h.services.UserService.Create(context.Background(), params)
}

func (h *Handler) getUserById(c *gin.Context) {

}

func (h *Handler) updateUser(c *gin.Context) {

}

func (h *Handler) deleteUser(c *gin.Context) {

}
