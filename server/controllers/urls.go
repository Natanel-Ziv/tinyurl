package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"tinyurl/server/models"
	"tinyurl/server/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
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

func (uc *URLController) GetLongURL(ctx *gin.Context) {
	short_hash := ctx.Param("hash")

	urlDetails, err := uc.urlService.GetURLFromShort(short_hash)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Err(err).Msg("failed to find short_url")
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "short url not found",
		})
		return
	}

	if time.Now().After(urlDetails.ExpiresAt) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "url expired",
		})
		return
	}

	err = uc.urlService.UpdateURLVisited(urlDetails.ID)
	if err != nil {
		log.Err(err).Msg("failed to update visited")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unexpected error",
		})
	}

	ctx.Redirect(http.StatusTemporaryRedirect, urlDetails.LongUrl)
}

func (uc *URLController) RegisterURL(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*models.UserDBResponse)

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

	urlReq.User = user.ID

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

func (uc *URLController) GetUserStatistics(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*models.UserDBResponse)

	allUserUrls, err := uc.urlService.GetAllURLForUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if allUserUrls == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "no urls for user",
		})
		return
	}

	var resp []*models.URLResponse
	for _, userUrl := range allUserUrls {
		resp = append(resp, models.URLFilteredResponse(userUrl))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"urls": resp,
		},
	})
}

func (uc *URLController) GetStatisticsForURL(ctx *gin.Context) {
	shortHash := ctx.Param("hash")
	user := ctx.MustGet("currentUser").(*models.UserDBResponse)

	urlDetails, err := uc.urlService.GetURLFromShort(shortHash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if urlDetails.User != user.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "this url does not belong to this user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"url": models.URLFilteredResponse(urlDetails),
		},
	})
}
