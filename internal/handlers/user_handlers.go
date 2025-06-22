package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"unicode"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/logger/sl"
	"github.com/gin-gonic/gin"
)

// getAllUsers godoc
// @Summary      get all users with optionally filters and pagination
// @Description  Get users using flexible query filters and pagination. You can provide partial values for `name`, `surname`, or `patronymic` — filtering will still work. Each of these parameters is optional and can be used independently or in combination.
// @Description
// @Description  Example: ?page=5&page_size=10
// @Description  Response: 10 users with offset=40
// @Description
// @Description  Example2: ?name=al
// @Description  Response: Alex, Alina, etc.
// @Description
// @Description  Example3: ?name=al&surname=sh
// @Description  Response: Alexandr Shprot, Alina Sham, etc.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        name        query     string  false "name filter"
// @Param        surname     query     string  false "surname filter"
// @Param        patronymic  query     string  false "patronymic filter"
// @Param        gender      query     string  false  "gender filter can be only male or female"
// @Param        page_size       query     int     false  "min:5"
// @Param        page      query     int     false  "min:1"
// @Success      200  {array}  entities.User
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	patronymic := c.DefaultQuery("patronymic", "")
	gender := c.DefaultQuery("gender", "")

	//validation
	pageSizeStr := c.DefaultQuery("page_size", "5")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 5 {
		pageSize = 5
		log.Debug("Invalid page_size value, set to 5", slog.String("user's page_size", pageSizeStr))
	}
	if pageSize > 50 {
		pageSize = 50
		log.Debug("page_size is greater than 50, set to 50", slog.String("user's page_size", pageSizeStr))
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
		log.Debug("Invalid page value, set to 1", slog.String("page", pageStr))
	}

	log.Info("Getting all users with parameters", slog.Int("page_size", pageSize), slog.Int("page", page), slog.String("name", name), slog.String("surname", surname), slog.String("patronymic", patronymic), slog.String("gender", gender))

	//I think using cache here might be useless because of variations of keys due to many filters
	allUsers, totalCount, err := h.services.UserService.GetAllUsers(pageSize, page, name, surname, patronymic, gender)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	if page == 1 && totalCount == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//check for invalid page num, totalPages also may be used by frontend but totalcount == 0 and status code 404 may be used here as well
	totalPages := math.Ceil(float64(totalCount) / float64(pageSize))
	if page > int(totalPages) {
		newErrorResponse(c, log, http.StatusBadRequest, "Page exceeds total number of pages", errors.New("page exceeds max of pages"))
		return
	}

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
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	var fullName entities.FullName
	err := c.ShouldBindJSON(&fullName)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Failed to parse json", err)
		return
	}

	if haveDigits(fullName.Name) || haveDigits(fullName.Surname) || haveDigits(fullName.Patronymic) {
		newErrorResponse(c, log, http.StatusBadRequest, "Name, Surname or Patronymic must not contain digits", errors.New("name,surname or patronimyc contain digits"))
		return
	}

	exists, err := h.services.UserService.ExistByFullName(fullName)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Failed to check if the user exists", err)
		return
	}
	if exists {
		newErrorResponse(c, log, http.StatusBadRequest, "User already exists", errors.New("user already exists"))
		return
	}

	age, gender, nationality, err := h.services.InfoRequestService.RequestAdditionalInfo(fullName.Name)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "Requests timed out or service is unreachable", err)
		return
	}

	params := entities.User{
		Name:        fullName.Name,
		Surname:     fullName.Surname,
		Patronymic:  fullName.Patronymic,
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

	log.Info("Created user successfully", slog.Any("created_user", createdUser))

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("User created successfully with id: %d", createdUser.Id)})
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
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	userIdStr := c.Param("user_id")
	userId32, err := parseInt32(userIdStr)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "Id should be number", err)
		return
	}

	//cache check
	var user entities.User
	cacheKey := "user:" + userIdStr
	err = h.services.RedisService.Get(context.Background(), cacheKey, &user)
	if err == nil {
		log.Info("User found in cache", slog.Int("user_id", int(userId32)))
		c.JSON(http.StatusOK, user)
		return
	}

	log.Info("Getting user by ID from postgres db", slog.Int("user_id", int(userId32)))
	user, err = h.services.UserService.GetUserById(userId32)
	if err != nil {
		newErrorResponse(c, log, http.StatusNotFound, "User not found", err)
		return
	}

	//cache set
	err = h.services.RedisService.Set(context.Background(), cacheKey, user)
	if err != nil {
		log.Warn("Failed to set user in cache", slog.Int("user_id", int(userId32)), sl.Err(err))
	}

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
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

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
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

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
