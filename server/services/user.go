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
	FindUserByID(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
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
func (us *UserServiceImpl) FindUserByID(id string) (*models.DBResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get objectId from hex: %w", err)
	}

	var user *models.DBResponse
	
	query := bson.M{"_id": oid}
	err = us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.DBResponse{}, err
		}
		return nil, fmt.Errorf("failed to find id: %w", err)
	}

	return user, nil

}

// FindUserByEmail implements UserService
func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	
	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.DBResponse{}, err
		}
		return nil, fmt.Errorf("failed to find id: %w", err)
	}

	return user, nil
}

