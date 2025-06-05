package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/logger/handlers/slogdiscard"
	service "github.com/Util787/junTask/internal/services"
	"github.com/Util787/junTask/internal/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createUser(t *testing.T) {

	tests := []struct {
		testname                string
		inputBody               string
		mockExistBehavior       func(s *mocks.MockUserService)
		mockInfoRequestBehavior func(s *mocks.MockInfoRequestService)
		mockCreateBehavior      func(s *mocks.MockUserService)
		expectedStatusCode      int
		expectedResponseBody    string
	}{
		// TODO: Add more test cases. Add time delta
		{
			testname:  "Ok",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(41, "female", "BY", nil)
			},
			mockCreateBehavior: func(s *mocks.MockUserService) {
				s.On("CreateUser", entities.User{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}).Return(entities.User{Id: 3, Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"message":"user created successfully with id: 3"}`,
		},
		{
			testname:                "Empty JSON",
			inputBody:               `{}`,
			mockExistBehavior:       func(s *mocks.MockUserService) {},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *mocks.MockUserService) {},
			expectedStatusCode:      400,
			expectedResponseBody:    `{"message":"Failed to parse json in createUser handler"}`,
		},
		{
			testname:  "Create service error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(41, "female", "BY", nil)
			},
			mockCreateBehavior: func(s *mocks.MockUserService) {
				s.On("CreateUser", entities.User{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}).Return(entities.User{}, errors.New("Something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Failed to create user"}`,
		},
		{
			testname:  "Exist service error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, errors.New("Something went wrong"))
			},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *mocks.MockUserService) {},
			expectedStatusCode:      500,
			expectedResponseBody:    `{"message":"Failed to check if the user exists"}`,
		},
		{
			testname:  "User already exists",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(true, nil)
			},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {},
			mockCreateBehavior:      func(s *mocks.MockUserService) {},
			expectedStatusCode:      400,
			expectedResponseBody:    `{"message":"User already exists"}`,
		},
		{
			testname:  "Info request error",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			mockInfoRequestBehavior: func(s *mocks.MockInfoRequestService) {
				s.On("RequestAdditionalInfo", "testname").Return(0, "", "", errors.New("api call unreachable"))
			},
			mockCreateBehavior:   func(s *mocks.MockUserService) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Requests timed out or service is unreachable"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := mocks.NewMockUserService(t)
			mockInfoRequestService := mocks.NewMockInfoRequestService(t)
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

func setupTestRouter(mockUserService *mocks.MockUserService, mockInfoRequestService *mocks.MockInfoRequestService, redisService *mocks.MockRedisService) *gin.Engine {
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
		mockGetAllUsersBehavior func(s *mocks.MockUserService)
		expectedStatusCode      int
		expectedResponseBody    string
	}{
		{
			testname: "Ok",
			queryStr: "",
			mockGetAllUsersBehavior: func(s *mocks.MockUserService) {
				s.On("GetAllUsers", 5, 1, "", "", "", "").Return([]entities.User{{Name: "Aleksey", Surname: "Ivanov"}}, 1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"name":"Aleksey","surname":"Ivanov"`,
		},
		//TODO: add more testcases
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mockUserService := mocks.NewMockUserService(t)
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
