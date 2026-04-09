package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PortofolioInput struct {
	Name string `json:"name" binding:"required"`
}

type PortofolioStatusInput struct {
	IsActive bool `json:"is_active"`
}

type Portofolio struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	IsActive  bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
	DeletedAt time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
