package routes

import (
	"tinyurl/server/controllers"

	"github.com/gin-gonic/gin"
)

func apiRoutes(rg *gin.RouterGroup) {
	apiRouter := rg.Group("/api")
	{
		apiRouter.GET("/ping", controllers.Ping)
	}
}
