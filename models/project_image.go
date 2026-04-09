package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectImage struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty"`
	ProjectID primitive.ObjectID `json:"project_id,omitempty" bson:"project_id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	IsActive  bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDeleted bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	DeletedAt time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type ProjectImageInput struct {
	Image     string `json:"image,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

type ProjectImageUpdateInput struct {
	Image     string `json:"image,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

type ProjectImageStatusInput struct {
	IsActive bool `json:"is_active,omitempty"`
}
