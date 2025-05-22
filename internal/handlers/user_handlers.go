package handlers

import (
	"context"
	"net/http"

	"github.com/Util787/junTask/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	allUsers, err := h.services.UserService.GetAll(context.Background())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, allUsers)
}

func (h *Handler) createUser(c *gin.Context) {
	var user entities.User
	err := c.BindJSON(&user)
	if err != nil {
		return
	}

	requestUserAdditionalInfo(c,user.Name)
}

func (h *Handler) getUserById(c *gin.Context) {

}

func (h *Handler) updateUser(c *gin.Context) {

}

func (h *Handler) deleteUser(c *gin.Context) {

}
