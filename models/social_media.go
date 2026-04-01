package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialMedia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Icon      string             `bson:"icon"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type SocialMediaCreateInput struct {
	Icon string `bson:"icon" binding:"required"`
	Name string `bson:"name" binding:"required"`
}

type SocialMediaResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Icon      string             `bson:"icon"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
