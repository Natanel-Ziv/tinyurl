package routes

import (
	"tinyurl/server/controllers"
	"tinyurl/server/middleware"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController *controllers.UserController
	userService    services.UserService
	authCfg        *controllers.AuthConfig
}

func NewUserRoutes(userController *controllers.UserController, userService services.UserService, authCfg *controllers.AuthConfig) *UserRoutes {
	return &UserRoutes{
		userController: userController,
		userService:    userService,
		authCfg:        authCfg,
	}
}

func (ur *UserRoutes) UserRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/users")

	router.Use(middleware.DeserializeUser(ur.userService, ur.authCfg))
	router.GET("/me", ur.userController.GetMe)
}
