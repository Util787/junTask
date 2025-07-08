package handlers

import (
	"log/slog"

	"github.com/Util787/user-manager-api/internal/logger/sl"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, log *slog.Logger, statusCode int, message string, err error) {
	log.Error(message, sl.Err(err))
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
