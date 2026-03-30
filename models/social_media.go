package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialMedia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Icon      string             `bson:"icon"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type SocialMediaCreateInput struct {
	Icon string `json:"icon" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type SocialMediaResponse struct {
	ID        string    `json:"id"`
	Icon      string    `json:"icon"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
