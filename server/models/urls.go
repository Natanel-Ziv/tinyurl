package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterURLInput struct {
	LongUrl   string    `json:"long_url" bson:"long_url" binding:"required,validateurl"`
	ShortHash string    `json:"short_hash" bson:"short_hash"`
	Visited   int       `json:"visited" bson:"visited"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type URLDBResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	LongUrl   string             `json:"long_url" bson:"long_url"`
	ShortHash string             `json:"short_hash" bson:"short_hash"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	Visited   int                `json:"visited" bson:"visited"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type URLResponse struct {
	LongUrl   string    `json:"long_url" bson:"long_url"`
	ShortHash string    `json:"short_hash" bson:"short_hash"`
	Visited   int       `json:"visited" bson:"visited"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func URLFilteredResponse(dbResp *URLDBResponse) *URLResponse {
	return &URLResponse{
		LongUrl:   dbResp.LongUrl,
		ShortHash: dbResp.ShortHash,
		ExpiresAt: dbResp.ExpiresAt,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
	}
}
