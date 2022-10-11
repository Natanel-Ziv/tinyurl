package controllers

import (
	"net/http"
	"strings"
	"time"
	"tinyurl/server/models"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
)

const DefaultExpTime = 14 * 24 * time.Hour // 14 Days

type URLController struct {
	urlService services.URLService
}

func NewURLController(urlService services.URLService) *URLController {
	return &URLController{
		urlService: urlService,
	}
}

func (uc *URLController) RegisterURL(ctx *gin.Context) {
	var urlReq *models.RegisterURLInput

	err := ctx.ShouldBindJSON(&urlReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if urlReq.ExpiresAt.IsZero() {
		urlReq.ExpiresAt = time.Now().Add(DefaultExpTime)
	}

	newUrl, err := uc.urlService.RegisterURL(urlReq)
	if err != nil {
		if strings.Contains(err.Error(), "short url already exist") {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		} else if strings.Contains(err.Error(), "generate") {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": models.URLFilteredResponse(newUrl),
		},
	})
}
