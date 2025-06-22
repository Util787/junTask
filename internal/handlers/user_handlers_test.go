package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/logger/handlers/slogdiscard"
	service "github.com/Util787/junTask/internal/services"
	serviceMock "github.com/Util787/junTask/internal/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_createUser(t *testing.T) {

	tests := []struct {
		testname                string
		inputBody               string
		mockExistBehavior       func(s *serviceMock.MockUserService)
		mockInfoRequestBehavior func(s *serviceMock.MockInfoRequestService)
		mockCreateBehavior      func(s *serviceMock.MockUserService)
		expectedStatusCode      int
		expectedResponseBody    string
	}{
		{
			testname:  "Ok",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(41, "female", "BY", nil)
			},
			mockCreateBehavior: func(s *serviceMock.MockUserService) {
				s.On("CreateUser", entities.User{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}).Return(entities.User{Id: 3, Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"message":"User created successfully with id: 3"}`,
		},
		{
			testname:                "Empty JSON",
			inputBody:               `{}`,
			mockExistBehavior:       func(s *serviceMock.MockUserService) {},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *serviceMock.MockUserService) {},
			expectedStatusCode:      400,
			expectedResponseBody:    `{"message":"Failed to parse json"}`,
		},
		{
			testname:  "Create service error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(41, "female", "BY", nil)
			},
			mockCreateBehavior: func(s *serviceMock.MockUserService) {
				s.On("CreateUser", entities.User{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}).Return(entities.User{}, errors.New("Something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Failed to create user"}`,
		},
		{
			testname:  "Exist service error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, errors.New("Something went wrong"))
			},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *serviceMock.MockUserService) {},
			expectedStatusCode:      500,
			expectedResponseBody:    `{"message":"Failed to check if the user exists"}`,
		},
		{
			testname:  "User already exists",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(true, nil)
			},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *serviceMock.MockUserService) {},
			expectedStatusCode:      400,
			expectedResponseBody:    `{"message":"User already exists"}`,
		},
		{
			testname:  "Info request error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *serviceMock.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(0, "", "", errors.New("api call unreachable"))
			},
			mockCreateBehavior:   func(s *serviceMock.MockUserService) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Requests timed out or service is unreachable"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := serviceMock.NewMockUserService(t)
			mockInfoRequestService := serviceMock.NewMockInfoRequestService(t)
			router := setupTestRouter(mockUserService, mockInfoRequestService, nil)

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			test.mockExistBehavior(mockUserService)
			test.mockInfoRequestBehavior(mockInfoRequestService)
			test.mockCreateBehavior(mockUserService)

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)
			assert.Contains(t, resp.Body.String(), test.expectedResponseBody)

		})
	}
}

func setupTestRouter(mockUserService *serviceMock.MockUserService, mockInfoRequestService *serviceMock.MockInfoRequestService, redisService *serviceMock.MockRedisService) *gin.Engine {
	logger := slogdiscard.NewDiscardLogger()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	service := &service.Service{UserService: mockUserService, InfoRequestService: mockInfoRequestService, RedisService: redisService}
	h := NewHandlers(service, logger)

	router.GET("/users", h.getAllUsers)
	router.POST("/users", h.createUser)
	router.GET("/users/:user_id", h.getUserById)
	router.PATCH("/users/:user_id", h.updateUser)
	router.DELETE("/users/:user_id", h.deleteUser)

	return router
}

// Signature: GetAllUsers(pageSize, page int, name, surname, patronymic, gender string) (users []entities.User, totalCount int,err error)
func TestHandler_getAllUsers(t *testing.T) {

	tests := []struct {
		testname                string
		queryStr                string
		mockGetAllUsersBehavior func(s *serviceMock.MockUserService)
		expectedStatusCode      int
		expectedResponseBody    string
	}{
		{
			testname: "Ok",
			queryStr: "",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return([]entities.User{{Name: "Aleksey", Surname: "Ivanov"}}, 1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"name":"Aleksey","surname":"Ivanov"`,
		},
		{
			testname: "Invalid page_size below minimum",
			queryStr: "?page_size=2",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return([]entities.User{{Name: "Aleksey", Surname: "Ivanov"}}, 1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"name":"Aleksey","surname":"Ivanov"`,
		},
		{
			testname: "Invalid page_size above maximum",
			queryStr: "?page_size=100",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 50, 1, "", "", "", "").Return([]entities.User{{Name: "Aleksey", Surname: "Ivanov"}}, 1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"name":"Aleksey","surname":"Ivanov"`,
		},
		{
			testname: "Invalid page number",
			queryStr: "?page=-1",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return([]entities.User{{Name: "Aleksey", Surname: "Ivanov"}}, 1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"name":"Aleksey","surname":"Ivanov"`,
		},
		{
			testname: "No users found",
			queryStr: "",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return([]entities.User{}, 0, nil)
			},
			expectedStatusCode:   404,
			expectedResponseBody: "",
		},
		{
			testname: "No users found with page >1",
			queryStr: "?page=4",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 4, "", "", "", "").Return([]entities.User{}, 0, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `"Page exceeds total number of pages"`,
		},
		{
			testname: "Page exceeds total number of pages",
			queryStr: "?page=3",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 3, "", "", "", "").Return([]entities.User{}, 5, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Page exceeds total number of pages",
		},
		{
			testname: "Internal server error",
			queryStr: "",
			mockGetAllUsersBehavior: func(s *serviceMock.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return(nil, 0, errors.New("DB error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to get users",
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := serviceMock.NewMockUserService(t)
			router := setupTestRouter(mockUserService, nil, nil)

			resp := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/users"+test.queryStr, nil)

			test.mockGetAllUsersBehavior(mockUserService)

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)
			assert.Contains(t, resp.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getUserById(t *testing.T) {
	tests := []struct {
		testname           string
		userId             string
		mockRedisGet       func(s *serviceMock.MockRedisService)
		mockUserServiceGet func(s *serviceMock.MockUserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			testname: "Ok from cache",
			userId:   "1",
			mockRedisGet: func(s *serviceMock.MockRedisService) {
				s.On("Get", mock.Anything, "user:1", mock.AnythingOfType("*entities.User")).Run(func(args mock.Arguments) {
					u := args.Get(2).(*entities.User)
					*u = entities.User{Id: 1, Name: "CachedUser"}
				}).Return(nil)
			},
			mockUserServiceGet: func(s *serviceMock.MockUserService) {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"id":1`,
		},
		{
			testname: "Ok from DB, cache miss and set cache",
			userId:   "2",
			mockRedisGet: func(s *serviceMock.MockRedisService) {
				s.On("Get", mock.Anything, "user:2", mock.AnythingOfType("*entities.User")).Return(errors.New("redis: nil"))
				s.On("Set", mock.Anything, "user:2", entities.User{Id: 2, Name: "DBUser"}).Return(nil)
			},
			mockUserServiceGet: func(s *serviceMock.MockUserService) {
				s.On("GetUserById", int32(2)).Return(entities.User{Id: 2, Name: "DBUser"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"id":2`,
		},
		{
			testname:           "Invalid user ID param",
			userId:             "abc",
			mockRedisGet:       func(s *serviceMock.MockRedisService) {},
			mockUserServiceGet: func(s *serviceMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"Id should be number"`,
		},
		{
			testname: "User not found in DB",
			userId:   "3",
			mockRedisGet: func(s *serviceMock.MockRedisService) {
				s.On("Get", mock.Anything, "user:3", mock.AnythingOfType("*entities.User")).Return(errors.New("redis: nil"))
			},
			mockUserServiceGet: func(s *serviceMock.MockUserService) {
				s.On("GetUserById", int32(3)).Return(entities.User{}, errors.New("not found"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `"User not found"`,
		},
		{
			testname: "Cache set warning ignored",
			userId:   "4",
			mockRedisGet: func(s *serviceMock.MockRedisService) {
				s.On("Get", mock.Anything, "user:4", mock.AnythingOfType("*entities.User")).Return(errors.New("redis: nil"))
				s.On("Set", mock.Anything, "user:4", entities.User{Id: 4, Name: "DBUser4"}).Return(errors.New("redis set error"))
			},
			mockUserServiceGet: func(s *serviceMock.MockUserService) {
				s.On("GetUserById", int32(4)).Return(entities.User{Id: 4, Name: "DBUser4"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"id":4`,
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := serviceMock.NewMockUserService(t)
			mockRedisService := serviceMock.NewMockRedisService(t)
			router := setupTestRouter(mockUserService, nil, mockRedisService)

			if test.mockRedisGet != nil {
				test.mockRedisGet(mockRedisService)
			}
			if test.mockUserServiceGet != nil {
				test.mockUserServiceGet(mockUserService)
			}

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users/"+test.userId, nil)

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)
			assert.Contains(t, resp.Body.String(), test.expectedResponse)
		})
	}
}

func TestHandler_updateUser(t *testing.T) {
	tests := []struct {
		testname           string
		userId             string
		inputBody          string
		mockExistBehavior  func(s *serviceMock.MockUserService)
		mockUpdateBehavior func(s *serviceMock.MockUserService)
		expectedCode       int
		expectedResponse   string
	}{
		{
			testname:           "Invalid user_id param",
			userId:             "abc",
			inputBody:          `{"name":"John"}`,
			mockExistBehavior:  func(s *serviceMock.MockUserService) {},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Id should be number",
		},
		{
			testname:  "ExistById service error",
			userId:    "1",
			inputBody: `{"name":"John"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(1)).Return(false, errors.New("db error"))
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Failed to check if the user exists",
		},
		{
			testname:  "User does not exist",
			userId:    "2",
			inputBody: `{"name":"John"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(2)).Return(false, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Cannot update user that does not exist",
		},
		{
			testname:  "Invalid JSON body",
			userId:    "3",
			inputBody: `{"name":123}`, // expecting string, not number
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(3)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Failed to parse json in updateUser handler",
		},
		{
			testname:  "Name contains digits",
			userId:    "4",
			inputBody: `{"name":"John123"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(4)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Invalid name",
		},
		{
			testname:  "Surname contains digits",
			userId:    "5",
			inputBody: `{"surname":"Smith1"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(5)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Invalid surname",
		},
		{
			testname:  "Patronymic contains digits",
			userId:    "6",
			inputBody: `{"patronymic":"Ivanov3"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(6)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Invalid patronymic",
		},
		{
			testname:  "Invalid gender",
			userId:    "7",
			inputBody: `{"gender":"unknown"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(7)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {},
			expectedCode:       http.StatusBadRequest,
			expectedResponse:   "Gender must be 'male' or 'female'",
		},
		{
			testname:  "UpdateUser service error",
			userId:    "8",
			inputBody: `{"name":"John"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(8)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {
				s.On("UpdateUser", int32(8), mock.Anything).Return(errors.New("update failed"))
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: "Failed to update user",
		},
		{
			testname:  "Success update",
			userId:    "9",
			inputBody: `{"name":"John","gender":"male"}`,
			mockExistBehavior: func(s *serviceMock.MockUserService) {
				s.On("ExistById", int32(9)).Return(true, nil)
			},
			mockUpdateBehavior: func(s *serviceMock.MockUserService) {
				s.On("UpdateUser", int32(9), mock.Anything).Return(nil)
			},
			expectedCode:     http.StatusOK,
			expectedResponse: "User updated successfully",
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := serviceMock.NewMockUserService(t)
			router := setupTestRouter(mockUserService, nil, nil)

			router.Use(func(c *gin.Context) {
				c.Next()
			})

			test.mockExistBehavior(mockUserService)
			test.mockUpdateBehavior(mockUserService)

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/users/"+test.userId, bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedCode, resp.Code)
			assert.Contains(t, resp.Body.String(), test.expectedResponse)
		})
	}
}
