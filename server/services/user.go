package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tinyurl/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	FindUserByID(string) (*models.UserDBResponse, error)
	FindUserByEmail(string) (*models.UserDBResponse, error)
}

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(ctx context.Context, collection *mongo.Collection) UserService {
	return &UserServiceImpl{
		collection: collection,
		ctx:        ctx,
	}
}

// FindUserByID implements UserService
func (us *UserServiceImpl) FindUserByID(id string) (*models.UserDBResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get objectId from hex: %w", err)
	}

	var user *models.UserDBResponse

	query := bson.M{"_id": oid}
	err = us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.UserDBResponse{}, err
		}
		return nil, fmt.Errorf("failed to find id: %w", err)
	}

	return user, nil

}

// FindUserByEmail implements UserService
func (us *UserServiceImpl) FindUserByEmail(email string) (*models.UserDBResponse, error) {
	var user *models.UserDBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.UserDBResponse{}, err
		}
		return nil, fmt.Errorf("failed to find id: %w", err)
	}

	return user, nil
}
