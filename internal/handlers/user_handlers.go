package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	allUsers, err := h.services.UserService.GetAll(context.Background())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, allUsers)
}

func (h *Handler) createUser(c *gin.Context) {
	// using entities.User instead of database.User is important to use binding tags
	var user entities.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in createUser handler")
		return
	}

	exists, err := h.services.UserService.Exist(context.Background(), database.UserExistsParams{
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: sql.NullString{String: user.Patronymic, Valid: true},
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error to check if the user exists")
		return
	}
	if exists {
		newErrorResponse(c, http.StatusBadRequest, "User already exists")
		return
	}

	age, gender, nationality := requestUserAdditionalInfo(c, user.Name)
	log.Println(user, age, gender, nationality)
	params := database.CreateUserParams{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        user.Name,
		Surname:     user.Surname,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
	if user.Patronymic == "" {
		params.Patronymic = sql.NullString{Valid: false}
	} else {
		params.Patronymic = sql.NullString{String: user.Patronymic, Valid: true}
	}

	createdUser, err := h.services.UserService.Create(context.Background(), params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Println("Created user: ", createdUser)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user created successfully with id: %v", createdUser.ID)})
}

func (h *Handler) getUserById(c *gin.Context) {

}

func (h *Handler) updateUser(c *gin.Context) {

}

func (h *Handler) deleteUser(c *gin.Context) {

}
