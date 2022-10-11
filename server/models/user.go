package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpInput struct {
	Email           string    `json:"email" bson:"email" binding:"required,validateemail"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UserDBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm" bson:"passwordConfirm,omitempty"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func UserFilteredResponse(dbResp *UserDBResponse) *UserResponse {
	return &UserResponse{
		ID:        dbResp.ID,
		Email:     dbResp.Email,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
	}
}
