package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Organizer   string             `bson:"organizer"`
	UrlFile     string             `bson:"url_file"`
	Year        string             `bson:"year"`
	Description string             `bson:"description,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type CreateCertificateInput struct {
	Title       string `form:"title" binding:"required"`
	Organizer   string `form:"organizer" binding:"required"`
	Year        string `form:"year" binding:"required"`
	Description string `form:"description"`
}

type UpdateCertificateInput struct {
	Title       *string `form:"title,omitempty"`
	Organizer   *string `form:"organizer,omitempty"`
	Year        *string `form:"year,omitempty"`
	Description *string `form:"description,omitempty"`
}
