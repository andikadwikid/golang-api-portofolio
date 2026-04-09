package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/database"
	"portofolio-api/models"
	"portofolio-api/utils"
)

// CreateSocialMedia godoc
// @Summary Create a new social media entry
// @Description Create a new social media entry with name and icon
// @Tags social_media
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param social_media body models.SocialMediaCreateInput true "Social media details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media [post]
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
	// userID, err := helpers.GetUserIDFromContext(c)
	// if err != nil {
	// 	c.JSON(401, gin.H{"error": err.Error()})
	// 	return
	// }

	// userCollection := database.DB.Collection("users")

	// var user models.User
	// // Check user at database
	// err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	// if err == mongo.ErrNoDocuments {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "User not found",
	// 	})
	// 	return
	// }

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// Mapping ke model
	newSocialMedia := models.SocialMedia{
		ID:        primitive.NewObjectID(),
		Icon:      input.Icon,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
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
		ID:        insertedID,
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

// GetSocialMedia godoc
// @Summary Get all social media entries
// @Description Retrieve a list of all active social media entries
// @Tags social_media
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media [get]
func GetSocialMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionSocialMedia := database.DB.Collection("social_media")
	cursor, err := collectionSocialMedia.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var socialMedias []models.SocialMediaResponse
	if err := cursor.All(ctx, &socialMedias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing social medias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Social medias fetched",
		"data":    socialMedias,
	})
}

func GetSocialMediaById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Validate ID
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	fmt.Println("ID:", id)

	// 2. Prepare collection
	collection := database.DB.Collection("social_media")

	// 3. Query
	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}

	var socialMedia models.SocialMediaResponse

	err = collection.FindOne(ctx, filter).Decode(&socialMedia)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Social media not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch data",
		})
		return
	}

	// 4. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media fetched successfully",
		"data":    socialMedia,
	})
}

// UpdateSocialMedia godoc
// @Summary Update a social media entry
// @Description Update social media details by ID
// @Tags social_media
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Social Media ID"
// @Param social_media body models.SocialMediaUpdateInput true "Updated social media details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media/{id} [put]
func UpdateSocialMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.DB.Collection("social_media")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.SocialMediaUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	// Normalize
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))

	update := bson.M{
		"$set": bson.M{
			"icon":       input.Icon,
			"name":       input.Name,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id, "is_deleted": false}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Social media already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Social media updated",
		"data":    input,
	})
}

// DeleteSocialMedia godoc
// @Summary Delete a social media entry
// @Description Soft delete a social media entry by ID
// @Tags social_media
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Social Media ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media/{id} [delete]
func DeleteSocialMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.DB.Collection("social_media")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id, "is_deleted": false}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media deleted successfully"})
}
