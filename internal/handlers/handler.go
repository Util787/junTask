package handlers

import (
	"log/slog"

	_ "github.com/Util787/user-manager-api/docs"
	"github.com/Util787/user-manager-api/internal/handlers/middleware"
	service "github.com/Util787/user-manager-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	log      *slog.Logger
	services *service.Service
}

func NewHandlers(services *service.Service, log *slog.Logger) *Handler {
	return &Handler{services: services, log: log}
}

func (h *Handler) InitRoutes(env string) *gin.Engine {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	if env != "prod" {
		router.Use(gin.Logger())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware(h.log))
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
