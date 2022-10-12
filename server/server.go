package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tinyurl/server/config"
	"tinyurl/server/controllers"
	"tinyurl/server/db"
	"tinyurl/server/routes"
	"tinyurl/server/services"
	"tinyurl/server/utils"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server      *http.Server
	ctx         context.Context
	cancelCtx   context.CancelFunc
	mongoClient *db.MongoDB
	redisClient *db.RedisDB
}

func NewServer(parentCtx context.Context, cfg *config.Config) (*Server, error) {
	ctx, cancelCtx := context.WithCancel(parentCtx)

	mongoClient, err := db.NewMongoDB(ctx, cfg.MongoDBUri, cfg.MongoDBName)
	if err != nil {
		return nil, err
	}

	redisClient, err := db.NewRedisDB(ctx, cfg.RedisUri)
	if err != nil {
		return nil, err
	}

	usersCollection := mongoClient.GetCollection("users")
	userService := services.NewUserServiceImpl(ctx, usersCollection)
	authService := services.NewAuthServiceImpl(ctx, usersCollection)

	urlsCollection := mongoClient.GetCollection("urls")
	urlService := services.NewURLServiceImpl(ctx, urlsCollection)

	authCfg := &controllers.AuthConfig{
		AccessTokenPrivateKey:  cfg.AccessTokenPrivateKey,
		AccessTokenPublicKey:   cfg.AccessTokenPublicKey,
		RefreshTokenPrivateKey: cfg.RefreshTokenPrivateKey,
		RefreshTokenPublicKey:  cfg.RefreshTokenPublicKey,
		AccessTokenExpiresIn:   cfg.AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  cfg.RefreshTokenExpiresIn,
		AccessTokenMaxAge:      cfg.AccessTokenMaxAge,
		RefreshTokenMaxAge:     cfg.RefreshTokenMaxAge,
	}

	authController := controllers.NewAuthController(authCfg, authService, userService)
	userController := controllers.NewUserController(userService)
	urlController := controllers.NewURLController(urlService)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	utils.InitValidators()

	router := gin.New()
	router.Use(ginzerolog.Logger("server"), cors.New(corsConfig))

	routes := routes.NewRoutes(authController, userService, userController, urlService, urlController)
	routes.InitRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	return &Server{
		server:      server,
		ctx:         ctx,
		cancelCtx:   cancelCtx,
		mongoClient: mongoClient,
		redisClient: redisClient,
	}, nil
}

func (server *Server)Start() error {
	err := server.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (server *Server) Close() error {
	defer server.cancelCtx()

	err := server.server.Shutdown(server.ctx)
	if err != nil {
		log.Err(err).Msg("Failed to close server")
	}

	err = server.mongoClient.Close()
	if err != nil {
		log.Err(err).Msg("Failed to close mongo")
	}

	err = server.redisClient.Close()
	if err != nil {
		log.Err(err).Msg("Failed to close redis")
	}

	return nil
}
