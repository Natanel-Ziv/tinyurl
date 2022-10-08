package routes

import (
	"tinyurl/server/controllers"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	authRoutes *AuthRoutes
	userRoutes *UserRoutes
}

func NewRoutes(authController *controllers.AuthController, userService services.UserService, userController *controllers.UserController) *Routes {
	return &Routes{
		authRoutes: NewAuthRoutes(authController, userService),
		userRoutes : NewUserRoutes(userController, userService, authController.GetAuthCfg()),
	}
}

func (routes *Routes) InitRoutes(router *gin.Engine) {
	routerGroup := router.Group("/v1")

	apiRoutes(routerGroup)
	routes.authRoutes.AuthRoutes(routerGroup)
	routes.userRoutes.UserRoutes(routerGroup)
}
