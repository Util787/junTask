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

// getAllUsers godoc
// @Summary      get all users with optionally filters and pagination
// @Description  Get users using flexible query filters and pagination. You can provide partial values for `name`, `surname`, or `patronymic` — filtering will still work. Each of these parameters is optional and can be used independently or in combination.
// @Description  Example: ?limit=5&offset=10
// @Description  Response: 5 users with offset=10
// @Description  Example: ?name=al
// @Description  Response: Alex, Alina, etc.
// @Description  Example2: ?name=al&surname=sh
// @Description  Response: Alexandr Shprot, Alina Sham, etc.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name        query     string  false "name filter"
// @Param        surname     query     string  false "surname filter"
// @Param        patronymic  query     string  false "patronymic filter"
// @Param        gender      query     string  false  "gender filter can be only male or female"
// @Param        limit       query     int     false  "default:0"
// @Param        offset      query     int     false  "default:0"
// @Success      200  {array}  entities.User
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [get]
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

	logrus.Infof("Requested to get all users with queries: limit:%d offset:%d name:%s surname:%s patronymic:%s gender:%s", parsedLimit, parsedOffset, name, surname, patronymic, gender)
	allUsers, err := h.services.UserService.GetAllUsers(int(parsedLimit), int(parsedOffset), name, surname, patronymic, gender)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get users", err)
		return
	}
	logrus.Debugf("Got %d users ", len(allUsers))

	c.JSON(http.StatusOK, allUsers)
}

// createUser godoc
// @Summary      create user
// @Description  creating new user with provided name, surname, patronymic(optional)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        fullname  body  entities.FullName  true  "Users fullname: name, surname, patronymic(optional)"
// @Success      201  {object}  map[string]string "message with created user's id"
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var user entities.FullName
	err := c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed to parse json in createUser handler: ", err)
		return
	}

	if haveDigits(user.Name) || haveDigits(user.Surname) || haveDigits(user.Patronymic) {
		newErrorResponse(c, http.StatusBadRequest, "Name, Surname or Patronymic must not contain digits", err)
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
	logrus.Debugf("Received additional info: age=%d, gender=%s, nationality=%s", age, gender, nationality)

	params := entities.User{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	logrus.Infof("Attempt to create user with parameters:%+v ", params)
	createdUser, err := h.services.UserService.CreateUser(params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "", err)
		return
	}
	logrus.Debugf("Created user: %+v", createdUser)

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("user created successfully with id: %v", createdUser.Id)})
}

// getUserById godoc
// @Summary      get user by id
// @Description  recieve user info by providing id in path
// @Tags         users
// @Produce      json
// @Param        user_id  path      int  true "user_id"
// @Success      200      {object}  entities.User
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /users/{user_id} [get]
func (h *Handler) getUserById(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Id should be number", err)
		return
	}

	logrus.Infof("Attempt to get user by ID=%d", userId32)
	user, err := h.services.UserService.GetUserById(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "User doesnt exist", err)
		return
	}
	logrus.Debugf("Got user:%v by Id=%d", user, userId32)

	c.JSON(http.StatusOK, user)
}

// updateUser godoc
// @Summary      update user info by id
// @Description  updating user info by id provided in path. In request body you can optionally provide: name, surname, patronymic, age, gender, nationality. Update_at will change automatically
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  path      int                     true "user_id"
// @Param        user     body      entities.UpdateUserParams  true "parameters for update"
// @Success      200      {object}  map[string]string       "message about user update"
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /users/{user_id} [patch]
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
		newErrorResponse(c, http.StatusBadRequest, "Name, Surname or Patronymic must not contain digits", err)
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

	logrus.Infof("Attempt to update user ID=%d with params: %+v", userId32, params)
	err = h.services.UserService.UpdateUser(params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	logrus.Debug("updated user: ", userIdStr)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// deleteUser godoc
// @Summary      delete user by id
// @Description  deleting user by id if exists
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  path      int  true "user_id"
// @Success      200      {object}  map[string]string  "successful deleting message"
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /users/{user_id} [delete]
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

	logrus.Infof("Attempt to delete user ID=%d", userId32)
	err = h.services.UserService.DeleteUser(userId32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant find user", err)
		return
	}
	logrus.Debug("Deleted user: ", userIdStr)

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
