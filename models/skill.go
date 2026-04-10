package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Skill struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	UserID    primitive.ObjectID `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type CreateSkillInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateSkillInput struct {
	Name *string `json:"name,omitempty"`
}
