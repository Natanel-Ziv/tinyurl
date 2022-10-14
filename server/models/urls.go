package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterURLInput struct {
	LongUrl   string             `json:"long_url" bson:"long_url" binding:"required,validateurl"`
	ShortHash string             `json:"short_hash" bson:"short_hash"`
	Visited   int                `json:"visited" bson:"visited"`
	User      primitive.ObjectID `json:"user" bson:"user" bindinig:"required"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type URLDBResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	LongUrl   string             `json:"long_url" bson:"long_url"`
	ShortHash string             `json:"short_hash" bson:"short_hash"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	Visited   int                `json:"visited" bson:"visited"`
	User      primitive.ObjectID `json:"user" bson:"user"`
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

type AllURLsResponse struct {
	UserID      primitive.ObjectID `json:"_id" bson:"_id"`
	AllUserURLs []URLDBResponse    `json:"all_user_urls" bson:"all_user_urls"`
}

func URLFilteredResponse(dbResp *URLDBResponse) *URLResponse {
	return &URLResponse{
		LongUrl:   dbResp.LongUrl,
		ShortHash: dbResp.ShortHash,
		Visited:   dbResp.Visited,
		ExpiresAt: dbResp.ExpiresAt,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
	}
}
