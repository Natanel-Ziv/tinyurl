package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tinyurl/server/models"
	"tinyurl/server/services"
	"tinyurl/server/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

const InvalidEmailOrPassword = "Invalid email or password"

type AuthConfig struct {
	AccessTokenPrivateKey  string
	AccessTokenPublicKey   string
	RefreshTokenPrivateKey string
	RefreshTokenPublicKey  string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	AccessTokenMaxAge      int
	RefreshTokenMaxAge     int
}

type AuthController struct {
	authConfig  *AuthConfig
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authConfig *AuthConfig, authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		authConfig:  authConfig,
		authService: authService,
		userService: userService,
	}
}

func (ac *AuthController) GetAuthCfg() *AuthConfig {
	return ac.authConfig
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Passwords do not match",
		})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": models.UserFilteredResponse(newUser),
		},
	})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	err := ctx.ShouldBindJSON(&credentials)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": InvalidEmailOrPassword,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	err = utils.VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": InvalidEmailOrPassword,
		})
	}

	// Generate Tokens
	access_token, err := utils.CreateToken(ac.authConfig.AccessTokenExpiresIn, user.ID, ac.authConfig.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	refresh_token, err := utils.CreateToken(ac.authConfig.RefreshTokenExpiresIn, user.ID, ac.authConfig.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", access_token, ac.authConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, ac.authConfig.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.authConfig.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"access_token": access_token,
	})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": message,
		})
		return
	}

	sub, err := utils.ValidateToken(cookie, ac.authConfig.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	user, err := ac.userService.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "the user belonging to this token no logger exists",
		})
		return
	}

	access_token, err := utils.CreateToken(ac.authConfig.AccessTokenExpiresIn, user.ID, ac.authConfig.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", access_token, ac.authConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.authConfig.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"access_token": access_token,
	})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
