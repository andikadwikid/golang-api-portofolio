package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLoginInput struct {
	ID       primitive.ObjectID `json:"id"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password" binding:"required,min=8"`
}

type CreateUserInput struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name" binding:"required,min=5,max=25"`
	Username  string             `json:"username" binding:"required,min=5,max=25"`
	Email     string             `json:"email" binding:"required,email"`
	Password  string             `json:"password" binding:"required,min=8"`
	Avatar    string             `bson:"avatar,omitempty" json:"avatar"`
	Bio       string             `bson:"bio,omitempty" json:"bio"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type UpdateUserInput struct {
	Name     *string `json:"name" binding:"omitempty,min=5,max=25"`
	Username *string `json:"username" binding:"omitempty,min=5,max=25"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=8"`
	Avatar   *string `json:"avatar"`
	Bio      *string `json:"bio"`
}

type UserResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}
