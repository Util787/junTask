package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode"

	"github.com/Util787/junTask/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	patronymic := c.DefaultQuery("patronymic", "")
	gender := c.DefaultQuery("gender", "")
	limitStr := c.DefaultQuery("limit", "0")
	parsedLimit, err := strconv.Atoi(limitStr)
	if err != nil {
		parsedLimit = 0
	}
	offsetStr := c.DefaultQuery("offset", "0")
	parsedOffset, err := strconv.Atoi(offsetStr)
	if err != nil {
		parsedOffset = 0
	}

	allUsers, err := h.services.UserService.GetAllUsers(int(parsedLimit), int(parsedOffset), name, surname, patronymic, gender)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "uncorrect query params", err)
	}
	c.JSON(http.StatusOK, allUsers)
}

func (h *Handler) createUser(c *gin.Context) {
	var user entities.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in createUser handler: ", err)
		return
	}

	if haveDigits(user.Name) || haveDigits(user.Surname) || haveDigits(user.Patronymic) {
		newErrorResponse(c, http.StatusBadRequest, "Name, Surname or Patronymic mustn contain digits", err)
		return
	}

	exists, err := h.services.UserService.ExistByFullName(entities.FullName{
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to check if the user exists", err)
		return
	}
	if exists {
		newErrorResponse(c, http.StatusBadRequest, "User already exists", err)
		return
	}

	age, gender, nationality, err := requestUserAdditionalInfo(user.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "requests time out or unreachable", err)
		return
	}

	params := entities.User{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdUser, err := h.services.UserService.CreateUser(params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", err)
		return
	}
	logrus.Println("Created user: ", createdUser)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user created successfully with id: %v", createdUser.Id)})
}

func (h *Handler) getUserById(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number", err)
		return
	}

	user, err := h.services.UserService.GetUserById(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "User doesnt exist", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number", err)
		return
	}

	exists, err := h.services.UserService.ExistById(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to check if the user exists", err)
		return
	}
	if !exists {
		newErrorResponse(c, http.StatusBadRequest, "Cant update user that doesnt exist", err)
		return
	}

	var user entities.UpdateUserParams
	err = c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in updateUser handler", err)
		return
	}

	if haveDigits(user.Name) || haveDigits(user.Surname) || haveDigits(user.Patronymic) {
		newErrorResponse(c, http.StatusBadRequest, "Name, Surname or Patronymic mustn contain digits", err)
		return
	}

	params := entities.UpdateUserParams{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
		Id:          userId32,
	}

	err = h.services.UserService.UpdateUser(params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number", err)
		return
	}

	exist, err := h.services.UserService.ExistById(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to check if the user exists", err)
		return
	}
	if !exist {
		newErrorResponse(c, http.StatusBadRequest, "Cant delete user that doesnt exist", err)
		return
	}

	err = h.services.UserService.DeleteUser(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user", err)
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

func haveDigits(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
