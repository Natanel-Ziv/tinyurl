package routes

import (
	"tinyurl/server/controllers"
	"tinyurl/server/middleware"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController *controllers.AuthController
	userService services.UserService
}

func NewAuthRoutes(authController *controllers.AuthController, userService services.UserService) *AuthRoutes {
	return &AuthRoutes{authController: authController, userService: userService}
}

func (ar *AuthRoutes) AuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", ar.authController.SignUpUser)
	router.POST("/login", ar.authController.SignInUser)
	router.GET("/refresh", ar.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(ar.userService, ar.authController.GetAuthCfg()), ar.authController.LogoutUser)
}
