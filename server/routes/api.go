package routes

import (
	"tinyurl/server/controllers"

	"github.com/gin-gonic/gin"
)

func apiRoutes(superRoute *gin.Engine) {
	apiRouter := superRoute.Group("/api")
	{
		apiRouter.GET("/ping", controllers.Ping)
	}
}
