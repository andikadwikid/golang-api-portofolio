package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	Link         string             `json:"link,omitempty" bson:"link,omitempty"`
	PortofolioID primitive.ObjectID `json:"portofolio_id,omitempty" bson:"portofolio_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	IsActive     bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsDeleted    bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	DeletedAt    time.Time          `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type ProjectInput struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Link         string `json:"link"`
	PortofolioID string `json:"portofolio_id" binding:"required"`
}

type ProjectUpdateInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Link        string `json:"link"`
}

type ProjectStatusInput struct {
	IsActive bool `json:"is_active"`
}
