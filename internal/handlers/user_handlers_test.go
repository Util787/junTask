package handlers

import (
	"bytes"
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
	mockUserService := mocks.NewMockUserService(t)
	router := setupRouterWithMock(mockUserService)

	tests := []struct {
		testname             string
		inputBody            string
		mockCreateBehavior   func(s *mocks.MockUserService)
		mockExistBehavior    func(s *mocks.MockUserService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		// TODO: Add more test cases.
		{
			testname:  "ok",
			inputBody: `{"name":"testname","surname":"testsurname","patronymic":"testpatronymic"}`,
			mockCreateBehavior: func(s *mocks.MockUserService) {
				//Можно было бы сделать чтобы age gender nationality сверялись с теми что в ответах апишек но будем честны при каждом тесте делать запросы это расточительство
				s.On("CreateUser", entities.User{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}).Return(entities.User{Id: 1, Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic", Age: 41, Gender: "female", Nationality: "BY"}, nil)
			},
			mockExistBehavior: func(s *mocks.MockUserService) {
				s.On("ExistByFullName", entities.FullName{Name: "testname", Surname: "testsurname", Patronymic: "testpatronymic"}).Return(false, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"message":"user created successfully with id: 1"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			test.mockExistBehavior(mockUserService)
			test.mockCreateBehavior(mockUserService)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			assert.Contains(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func setupRouterWithMock(mockUserService *mocks.MockUserService) *gin.Engine {
	logger := slogdiscard.NewDiscardLogger()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	service := &service.Service{UserService: mockUserService}
	h := NewHandlers(service, logger)

	router.GET("/users", h.getAllUsers)
	router.POST("/users", h.createUser)
	router.GET("/users/:user_id", h.getUserById)
	router.PATCH("/users/:user_id", h.updateUser)
	router.DELETE("/users/:user_id", h.deleteUser)

	return router
}
