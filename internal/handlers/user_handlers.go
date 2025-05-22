package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")
	patronymic := c.Query("patronymic")
	gender := c.Query("gender")
	limit, err := parseInt32(c.Query("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Limit must be a number")
		return
	}
	offset, err := parseInt32(c.Query("offset"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Offset must be a number")
		return
	}

	allUsers, err := h.services.UserService.GetAllUsers(context.Background(), database.GetAllUsersParams{
		Limit:      limit,
		Offset:     offset,
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Gender:     gender,
	})
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

	exists, err := h.services.UserService.ExistByFullName(context.Background(), database.UserExistByFullNameParams{
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: sql.NullString{String: user.Patronymic, Valid: true},
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to check if the user exists")
		return
	}
	if exists {
		newErrorResponse(c, http.StatusBadRequest, "User already exists")
		return
	}

	age, gender, nationality, err := requestUserAdditionalInfo(user.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "requests time out or unreachable")
		return
	}

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

	// you can route this structure parameters to entities.User struct if you want to control json tags
	createdUser, err := h.services.UserService.CreateUser(context.Background(), params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Println("Created user: ", createdUser)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user created successfully with id: %v", createdUser.ID)})
}

func (h *Handler) getUserById(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	// you can route this structure parameters to entities.User struct if you want to control json tags
	user, err := h.services.UserService.GetUserById(context.Background(), userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	userBeforeUpdate, err := h.services.UserService.GetUserById(context.Background(), userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user")
		return
	}

	var user entities.UpdateUser
	err = c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in updateUser handler")
		return
	}

	// if json parameter in body is empty sqlc will put zerovalue in column, couldnt find any way to fix it so the only way is to check it manualy
	name := user.Name
	if name == "" {
		name = userBeforeUpdate.Name
	}

	surname := user.Surname
	if surname == "" {
		surname = userBeforeUpdate.Surname
	}

	age := user.Age
	if age == 0 {
		age = userBeforeUpdate.Age
	}

	gender := user.Gender
	if gender == "" {
		gender = userBeforeUpdate.Gender
	}

	nationality := user.Nationality
	if nationality == "" {
		nationality = userBeforeUpdate.Nationality
	}

	var patronymic sql.NullString
	if user.Patronymic == "" {
		patronymic = sql.NullString{String: userBeforeUpdate.Patronymic.String, Valid: userBeforeUpdate.Patronymic.Valid}
	} else {
		patronymic = sql.NullString{String: user.Patronymic, Valid: true}
	}

	params := database.UpdateUserParams{
		UpdatedAt:   time.Now(),
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
		ID:          userId32,
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

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number")
		return
	}

	exist, err := h.services.UserService.ExistById(context.Background(), userId32)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to check if the user exists")
		return
	}
	if !exist {
		newErrorResponse(c, http.StatusBadRequest, "Cant delete user that doesnt exist")
		return
	}

	err = h.services.UserService.DeleteUser(context.Background(), userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User with id:%s deleted successfully", userIdStr)})
}

func parseInt32(numStr string) (int32, error) {
	parsedNum, err := strconv.ParseInt(numStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsedNum), nil
}
