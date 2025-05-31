package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/Util787/junTask/entities"
	"github.com/gin-gonic/gin"
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
	//can add uuid to every operation individually to track it
	const op = "getAllUsers"
	log := h.log.With(
		slog.String("op", op),
	)
	start := time.Now()
	log.Debug("Start", slog.Time("start_time", start))

	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	patronymic := c.DefaultQuery("patronymic", "")
	gender := c.DefaultQuery("gender", "")

	limitStr := c.DefaultQuery("limit", "0")
	parsedLimit, err := strconv.Atoi(limitStr)
	if err != nil {
		parsedLimit = 0
		log.Debug("Invalid limit value, set to 0", slog.String("limit", limitStr))
	}

	offsetStr := c.DefaultQuery("offset", "0")
	parsedOffset, err := strconv.Atoi(offsetStr)
	if err != nil {
		parsedOffset = 0
		log.Debug("Invalid offset value, set to 0", slog.String("offset", offsetStr))
	}

	log.Info("Getting all users with parameters",
		slog.Int("limit", parsedLimit),
		slog.Int("offset", parsedOffset),
		slog.String("name", name),
		slog.String("surname", surname),
		slog.String("patronymic", patronymic),
		slog.String("gender", gender),
	)

	allUsers, err := h.services.UserService.GetAllUsers(int(parsedLimit), int(parsedOffset), name, surname, patronymic, gender)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	log.Debug("Operation finished", slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	log.Info("Got users successfully", slog.Int("count", len(allUsers)))

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
	const op = "createUser"
	log := h.log.With(
		slog.String("op", op),
	)
	start := time.Now()
	log.Debug("Start", slog.Time("start_time", start))

	var user entities.FullName
	err := c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Failed to parse json in createUser handler ", err)
		return
	}

	//validation
	if haveDigits(user.Name) || haveDigits(user.Surname) || haveDigits(user.Patronymic) {
		newErrorResponse(c, log, http.StatusBadRequest, "Name, Surname or Patronymic must not contain digits", errors.New("name,surname or patronimyc contain digits"))
		return
	}

	exists, err := h.services.UserService.ExistByFullName(user)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to check if the user exists", err)
		return
	}
	if exists {
		newErrorResponse(c, log, http.StatusBadRequest, "User already exists", errors.New("user already exists"))
		return
	}

	log.Info("Requesting additional info")
	age, gender, nationality, err := requestUserAdditionalInfo(user.Name)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Requests timed out or service is unreachable", err)
		return
	}

	log.Debug("Received additional info",
		slog.Int("age", age),
		slog.String("gender", gender),
		slog.String("nationality", nationality),
	)

	params := entities.User{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	log.Info("Creating user with parameters", slog.Any("user", params))
	createdUser, err := h.services.UserService.CreateUser(params)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	log.Debug("Operation finished", slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	log.Info("Created user successfully", slog.Any("created_user", createdUser))

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
	const op = "getUserById"
	log := h.log.With(
		slog.String("op", op),
	)
	start := time.Now()
	log.Debug("Start", slog.Time("start_time", start))

	userIdStr := c.Param("user_id")
	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Id should be number", err)
		return
	}

	log.Info("Getting user by ID", slog.Int("user_id", int(userId32)))
	user, err := h.services.UserService.GetUserById(userId32)
	if err != nil {
		newErrorResponse(c, log, http.StatusNotFound, "User not found", err)
		return
	}

	log.Debug("Operation finished", slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	log.Info("Got user successfully", slog.Any("user", user))

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
	const op = "updateUser"
	log := h.log.With(
		slog.String("op", op),
	)
	start := time.Now()
	log.Debug("Start", slog.Time("start_time", start))

	userIdStr := c.Param("user_id")

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Id should be number", err)
		return
	}

	log.Info("Checking if user exists by id", slog.Int("user_id", int(userId32)))
	exists, err := h.services.UserService.ExistById(userId32)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Failed to check if the user exists", err)
		return
	}
	if !exists {
		newErrorResponse(c, log, http.StatusBadRequest, "Cannot update user that does not exist", errors.New("user does not exist"))
		return
	}
	log.Info("User exists in db", slog.Int("user_id", int(userId32)))

	var user entities.UpdateUserParams
	err = c.ShouldBindJSON(&user)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Failed to parse json in updateUser handler", err)
		return
	}

	//validation
	if user.Name != nil && haveDigits(*user.Name) {
		newErrorResponse(c, log, http.StatusBadRequest, "Name must not contain digits", errors.New("name contain digits"))
		return
	}
	if user.Surname != nil && haveDigits(*user.Surname) {
		newErrorResponse(c, log, http.StatusBadRequest, "Surname must not contain digits", errors.New("surname contain digits"))
		return
	}
	if user.Patronymic != nil && haveDigits(*user.Patronymic) {
		newErrorResponse(c, log, http.StatusBadRequest, "Patronymic must not contain digits", errors.New("patronimyc contain digits"))
		return
	}
	if user.Gender != nil && *user.Gender != "female" && *user.Gender != "male" {
		newErrorResponse(c, log, http.StatusBadRequest, "Gender must be 'male' or 'female'", errors.New("unknown gender"))
		return
	}

	log.Info("Updating user with parameters", slog.Any("update_params", user))
	err = h.services.UserService.UpdateUser(userId32, user)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	log.Debug("Operation finished", slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	log.Info("Updated user successfully", slog.Int("user_id", int(userId32)))

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
	const op = "deleteUser"
	log := h.log.With(
		slog.String("op", op),
	)
	start := time.Now()
	log.Debug("Start", slog.Time("start_time", start))

	userIdStr := c.Param("user_id")

	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Id should be number", err)
		return
	}

	log.Info("Checking if user exists by id", slog.Int("user_id", int(userId32)))
	exist, err := h.services.UserService.ExistById(userId32)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to check if the user exists", err)
		return
	}
	if !exist {
		newErrorResponse(c, log, http.StatusBadRequest, "Cannot delete user that does not exist", errors.New("user does not exist"))
		return
	}
	log.Info("User exists in db", slog.Int("user_id", int(userId32)))

	log.Info("Deleting user", slog.Int("user_id", int(userId32)))
	err = h.services.UserService.DeleteUser(userId32)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Cannot find user", err)
		return
	}
	log.Debug("Operation finished", slog.Int64("duration_ms", time.Since(start).Milliseconds()))
	log.Info("Deleted user successfully", slog.Int("user_id", int(userId32)))

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
