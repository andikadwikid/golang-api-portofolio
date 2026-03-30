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

	"portofolio-api/database"
	"portofolio-api/helpers"
	"portofolio-api/models"
	"portofolio-api/utils"
)

func CreateSocialMedia(c *gin.Context) {
	var input models.SocialMediaCreateInput

	// Bind input
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	// Normalize
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.DB.Collection("social_media")

	// Get user ID from context
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.DB.Collection("users")

	var user models.User
	// Check user at database
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Mapping ke model
	newSocialMedia := models.SocialMedia{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Icon:      input.Icon,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, newSocialMedia)
	if err != nil {

		// Handle duplicate key (dari unique index)
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Social media already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Safe casting
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ID"})
		return
	}

	response := models.SocialMediaResponse{
		ID:        insertedID.Hex(),
		Name:      newSocialMedia.Name,
		Icon:      newSocialMedia.Icon,
		CreatedAt: newSocialMedia.CreatedAt,
		UpdatedAt: newSocialMedia.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Social media created",
		"data":    response,
	})
}
