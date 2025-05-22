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
		lists := api.Group("/users")
		{
			lists.GET("/", h.getAllUsers)
			lists.POST("/", h.createUser)
			lists.GET("/:user_id", h.getUserById)
			lists.PUT("/:user_id", h.updateUser)
			lists.DELETE("/:user_id", h.deleteUser)
		}
	}
	return router
}
