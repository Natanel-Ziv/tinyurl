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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLService interface {
	RegisterURL(*models.RegisterURLInput) (*models.URLDBResponse, error)
	GetURLFromShort(string) (*models.URLDBResponse, error)
	UpdateURLVisited(primitive.ObjectID) error
	GetAllURLForUser(primitive.ObjectID) ([]*models.URLDBResponse, error)
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

// UpdateURLVisited implements URLService
func (us *URLServiceImpl) UpdateURLVisited(id primitive.ObjectID) error {
	query := bson.M{"_id": id}
	_, err := us.collection.UpdateOne(us.ctx, query, bson.D{{Key: "$inc", Value: bson.D{{Key: "visited", Value: 1}}}}, options.Update().SetUpsert(true))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		return fmt.Errorf("failed to find id: %w", err)
	}
	return nil
}

// RegisterURL implements URLService
func (us *URLServiceImpl) RegisterURL(urlRegisterRequest *models.RegisterURLInput) (*models.URLDBResponse, error) {
	urlRegisterRequest.CreatedAt = time.Now()
	urlRegisterRequest.UpdatedAt = urlRegisterRequest.CreatedAt
	urlRegisterRequest.Visited = 0

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
	index = mongo.IndexModel{Keys: bson.M{"long_url": 1}, Options: options.Index().SetUnique(true)}
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

func (us *URLServiceImpl) GetAllURLForUser(userId primitive.ObjectID) ([]*models.URLDBResponse, error) {
	query := bson.M{"user": userId}

	cur, err := us.collection.Find(us.ctx, query, &options.FindOptions{})
	if err != nil {
		return nil, err
	}

	var allUserUrls []*models.URLDBResponse
	if err = cur.All(us.ctx, &allUserUrls); err != nil {
		return nil, err
	}
	return allUserUrls, nil
}
