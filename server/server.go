package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tinyurl/server/config"
	"tinyurl/server/db"
	"tinyurl/server/routes"

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

	mongoClient, err := db.NewMongoDB(ctx, cfg.MongoDBUri)
	if err != nil {
		return nil, err
	}

	redisClient, err := db.NewRedisDB(ctx, cfg.RedisUri)
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	routes.AddRoutes(router)

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

func (server *Server) Start() error {
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
