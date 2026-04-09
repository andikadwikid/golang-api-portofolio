package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialMediaUser struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Link          string             `bson:"link"`
	SocialMediaID primitive.ObjectID `bson:"social_media_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	IsActive      bool               `bson:"is_active"`
}

type SocialMediaUserInput struct {
	Link          string             `json:"link" binding:"required"`
	SocialMediaID primitive.ObjectID `json:"social_media_id" binding:"required"`
}

type SocialMediaUserStatusUpdateInput struct {
	IsActive bool `json:"is_active"`
}

type SocialMediaUserUpdateInput struct {
	Link          string             `json:"link" binding:"required"`
	SocialMediaID primitive.ObjectID `json:"social_media_id" binding:"required"`
}

type SocialMediaUserResponse struct {
	ID            primitive.ObjectID `json:"id"`
	Link          string             `json:"link"`
	SocialMediaID primitive.ObjectID `json:"social_media_id"`
	UserID        primitive.ObjectID `json:"user_id"`
	IsActive      bool               `json:"is_active"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}
