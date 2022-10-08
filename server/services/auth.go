package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tinyurl/server/models"
	"tinyurl/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignInInput) (*models.DBResponse, error)
}

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthServiceImpl(ctx context.Context, collection *mongo.Collection) AuthService {
	return &AuthServiceImpl{collection: collection, ctx: ctx}
}

// SignInUser implements AuthService
func (as *AuthServiceImpl) SignInUser(user *models.SignInInput) (*models.DBResponse, error) {
	// TODO: use me!
	return nil, nil
}

// SignUpUser implements AuthService
func (as *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password %w", err)
	}
	user.Password = hashedPassword

	res, err := as.collection.InsertOne(as.ctx, &user)
	if err != nil {
		if innerErr, ok := err.(mongo.WriteException); ok && innerErr.WriteErrors[0].HasErrorCode(11000) {
			return nil, errors.New("user with such mail already exists")
		}
		return nil, err
	}


	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}

	_, err = as.collection.Indexes().CreateOne(as.ctx, index)
	if err != nil {
		return nil, fmt.Errorf("could not create index for email %w", err)
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = as.collection.FindOne(as.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil

}

