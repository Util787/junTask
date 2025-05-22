package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	userIdStr := c.Param("user_id")
	userID64, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	user, err := h.services.UserService.GetUserById(context.Background(), int32(userID64))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userID64, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	var user entities.UpdateUser
	err = c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in updateUser handler")
		return
	}

	params := database.UpdateUserParams{
		UpdatedAt:   time.Now(),
		Name:        user.Name,
		Surname:     user.Surname,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
		ID:          int32(userID64),
	}
	if user.Patronymic == "" {
		params.Patronymic = sql.NullString{Valid: false}
	} else {
		params.Patronymic = sql.NullString{String: user.Patronymic, Valid: true}
	}

	err = h.services.UserService.UpdateUser(context.Background(), params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userID64, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	err = h.services.UserService.DeleteUser(context.Background(), int32(userID64))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User with id:%s deleted successfully", userIdStr)})
}
