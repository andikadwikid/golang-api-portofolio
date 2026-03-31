package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Avatar    string             `bson:"avatar,omitempty"`
	Bio       string             `bson:"bio,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserLoginInput struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email" binding:"required,email"`
	Password string             `bson:"password" binding:"required,min=8"`
}

type CreateUserInput struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" binding:"required,min=5,max=25"`
	Username  string             `bson:"username" binding:"required,min=5,max=25"`
	Email     string             `bson:"email" binding:"required,email"`
	Password  string             `bson:"password" binding:"required,min=8"`
	Avatar    string             `bson:"avatar,omitempty"`
	Bio       string             `bson:"bio,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type UpdateUserInput struct {
	Name     *string `bson:"name" binding:"omitempty,min=5,max=25"`
	Username *string `bson:"username" binding:"omitempty,min=5,max=25"`
	Email    *string `bson:"email" binding:"omitempty,email"`
	Password *string `bson:"password" binding:"omitempty,min=8"`
	Avatar   *string `bson:"avatar"`
	Bio      *string `bson:"bio"`
}

type UserResponse struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
