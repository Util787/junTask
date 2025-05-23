package handlers

import (
	service "github.com/Util787/junTask/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandlers(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
			users.POST("/", h.createUser)
			users.GET("/:user_id", h.getUserById)
			users.PATCH("/:user_id", h.updateUser)
			users.DELETE("/:user_id", h.deleteUser)
		}
	}
	return router
}
