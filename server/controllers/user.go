package controllers

import (
	"net/http"
	"tinyurl/server/models"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.UserDBResponse)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": models.UserFilteredResponse(currentUser),
		},
	})
}
