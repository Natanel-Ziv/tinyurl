package routes

import (
	"tinyurl/server/controllers"
	"tinyurl/server/middleware"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

type URLRoutes struct {
	urlController *controllers.URLController
	urlService    services.URLService
	userService   services.UserService
	authCfg       *controllers.AuthConfig
}

func NewURLRoutes(urlController *controllers.URLController, urlService services.URLService, userService services.UserService, authCfg *controllers.AuthConfig) *URLRoutes {
	return &URLRoutes{
		urlController: urlController,
		urlService:    urlService,
		userService:   userService,
		authCfg:       authCfg,
	}
}

func (ur *URLRoutes) URLRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/urls")

	router.Use(middleware.DeserializeUser(ur.userService, ur.authCfg))
	router.POST("/register", ur.urlController.RegisterURL)

	ur.statisticsRoutes(router)
}

func (ur *URLRoutes) statisticsRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/statistics")

	router.GET("/", ur.urlController.GetUserStatistics)
	router.GET("/:hash", ur.urlController.GetStatisticsForURL)
}

func (ur *URLRoutes) RedirectURLRoute(rg *gin.RouterGroup) {
	rg.GET(":hash", ur.urlController.GetLongURL)
}
