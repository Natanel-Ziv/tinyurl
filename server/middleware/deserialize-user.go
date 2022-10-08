package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"tinyurl/server/controllers"
	"tinyurl/server/services"
	"tinyurl/server/utils"

	"github.com/gin-gonic/gin"
)

func DeserializeUser(userService services.UserService, authCfg *controllers.AuthConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"message": "You are not logged in",
			})
			return
		}

		sub, err := utils.ValidateToken(access_token, authCfg.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"message": err.Error(),
			})
			return
		}

		user, err := userService.FindUserByID(fmt.Sprint(sub))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"message": "The user belonging to this token no logger exists",
			})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}