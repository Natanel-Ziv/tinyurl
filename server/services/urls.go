package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
	"tinyurl/server/models"
	"tinyurl/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLService interface {
	RegisterURL(*models.RegisterURLInput) (*models.URLDBResponse, error)
	GetURLFromShort(string) (*models.URLDBResponse, error)
}

type URLServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewURLServiceImpl(ctx context.Context, collection *mongo.Collection) URLService {
	return &URLServiceImpl{
		collection: collection,
		ctx:        ctx,
	}
}

// GetURLFromShort implements URLService
func (us *URLServiceImpl) GetURLFromShort(shortHash string) (*models.URLDBResponse, error) {
	var urlDetails *models.URLDBResponse

	query := bson.M{"short_hash": shortHash}
	err := us.collection.FindOne(us.ctx, query).Decode(&urlDetails)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.URLDBResponse{}, err
		}
		return nil, fmt.Errorf("failed to find id: %w", err)
	}

	return urlDetails, nil
}

// RegisterURL implements URLService
func (us *URLServiceImpl) RegisterURL(urlRegisterRequest *models.RegisterURLInput) (*models.URLDBResponse, error) {
	urlRegisterRequest.CreatedAt = time.Now()
	urlRegisterRequest.UpdatedAt = urlRegisterRequest.CreatedAt

	if urlRegisterRequest.ShortHash != "" {
		_, err := us.GetURLFromShort(urlRegisterRequest.ShortHash)
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("short hash already exist")
		}
	} else {
		shortURLSize, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random size: %w", err)
		}
		urlRegisterRequest.ShortHash, err = utils.GenerateRandomShortStringInSIze(int(shortURLSize.Int64()))
		if err != nil {
			return nil, fmt.Errorf("failed to generate short url: %w", err)
		}
	}
	
	res, err := us.collection.InsertOne(us.ctx, &urlRegisterRequest)
	if err != nil {
		if innerErr, ok := err.(mongo.WriteException); ok && innerErr.WriteErrors[0].HasErrorCode(11000) {
			return nil, errors.New("url with such long url already exists")
		}
		return nil, err
	}

	index := mongo.IndexModel{Keys: bson.M{"short_hash": 1}, Options: options.Index().SetUnique(true)}

	_, err = us.collection.Indexes().CreateOne(us.ctx, index)
	if err != nil {
		return nil, fmt.Errorf("could not create index for short_url %w", err)
	}

	var newShortURL *models.URLDBResponse
	query := bson.M{"_id": res.InsertedID}

	err = us.collection.FindOne(us.ctx, query).Decode(&newShortURL)
	if err != nil {
		return nil, err
	}

	return newShortURL, nil
}
