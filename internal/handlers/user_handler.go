package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	h.services.UserService.GetAll(context.Background())

	
}

func (h *Handler) createUser(c *gin.Context) {

}

func (h *Handler) getUserById(c *gin.Context) {

}

func (h *Handler) updateUser(c *gin.Context) {

}

func (h *Handler) deleteUser(c *gin.Context) {

}
