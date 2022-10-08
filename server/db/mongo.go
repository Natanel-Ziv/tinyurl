package db

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	ctx           context.Context
	mongoConn     *options.ClientOptions
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
}

func NewMongoDB(ctx context.Context, uri string, dbName string) (*MongoDB, error) {
	mongoConn := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Err(err).Msg("Error connecting mongo")
		return nil, err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Err(err).Msg("Error pinging mongo")
		return nil, err
	}

	log.Debug().Msg("Mongo connected!")

	return &MongoDB{
		ctx:           ctx,
		mongoConn:     mongoConn,
		mongoClient:   mongoClient,
		mongoDatabase: mongoClient.Database(dbName),
	}, nil
}

func (mogonDB *MongoDB) Close() error {
	return mogonDB.mongoClient.Disconnect(mogonDB.ctx)
}

func (mongoDB *MongoDB) GetCollection(name string) *mongo.Collection {
	return mongoDB.mongoDatabase.Collection(name)
}