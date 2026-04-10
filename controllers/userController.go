package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/config"
	"portofolio-api/database"
	"portofolio-api/models"
	"portofolio-api/repositories"
	"portofolio-api/utils"
)

type UserController struct {
	userRepo repositories.UserRepository
}

func NewUserController(userRepo repositories.UserRepository) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

// RegisterUser godoc
func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var input models.CreateUserInput

	// Bind & validate
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	// Normalize email
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check existing email
	_, err := ctrl.userRepo.FindByEmail(ctx, input.Email)

	if err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email already registered", nil)
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", nil)
		return
	}

	// Mapping ke model DB
	user := models.User{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
		Avatar:    "",
		Bio:       "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert ke Mongo via Repo
	err = ctrl.userRepo.Create(ctx, &user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"id": user.ID.Hex(),
	})
}

// LoginUser godoc
func (ctrl *UserController) LoginUser(c *gin.Context) {
	var input models.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email via Repo
	user, err := ctrl.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Check password
	if !utils.CheckPassword(input.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Generate JWT via Config
	token, err := utils.GenerateToken(user.ID.Hex(), string(config.Config.JWTSecret))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetUsers godoc
func (ctrl *UserController) GetUsers(c *gin.Context) {
	collection := database.DB.Collection(config.CollectionUsers)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer cursor.Close(ctx)

	var users []models.UserResponse
	if err := cursor.All(ctx, &users); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users fetched successfully", users)
}

// UpdateUser godoc
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	update := bson.M{}
	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Username != nil {
		update["username"] = *input.Username
	}
	if input.Email != nil {
		update["email"] = *input.Email
	}
	if input.Avatar != nil {
		update["avatar"] = *input.Avatar
	}
	if input.Bio != nil {
		update["bio"] = *input.Bio
	}
	update["updated_at"] = time.Now()

	collection := database.DB.Collection(config.CollectionUsers)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser godoc
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	collection := database.DB.Collection(config.CollectionUsers)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}
