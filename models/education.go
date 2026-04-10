package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Education struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	StartDate   time.Time          `bson:"start_date"`
	EndDate     time.Time          `bson:"end_date"`
	Score       string             `bson:"score,omitempty"`
	Description string             `bson:"description,omitempty"`
	Faculty     string             `bson:"faculty,omitempty"`
	Degree      string             `bson:"degree,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type CreateEducationInput struct {
	Name        string             `json:"name" binding:"required"`
	StartDate   time.Time          `json:"start_date" binding:"required"`
	EndDate     time.Time          `json:"end_date" binding:"required"`
	Score       string             `json:"score,omitempty"`
	Description string             `json:"description,omitempty"`
	Faculty     string             `json:"faculty,omitempty"`
	Degree      string             `json:"degree,omitempty"`
	UserID      primitive.ObjectID `json:"user_id"`
}

type UpdateEducationInput struct {
	Name        *string    `json:"name,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Score       *string    `json:"score,omitempty"`
	Description *string    `json:"description,omitempty"`
	Faculty     *string    `json:"faculty,omitempty"`
	Degree      *string    `json:"degree,omitempty"`
}
