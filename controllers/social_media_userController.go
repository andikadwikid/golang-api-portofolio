package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/database"
	"portofolio-api/helpers"
	"portofolio-api/models"
)

func CreateSocialMediaUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionSocialMediaUser := database.DB.Collection("social_media_user")

	collectionSocialMedia := database.DB.Collection("social_media")

	var input models.SocialMediaUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	socialMedia := collectionSocialMedia.FindOne(ctx, bson.M{"_id": input.SocialMediaID})
	if socialMedia == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Social media not found"})
		return
	}

	newSocialMediaUser := models.SocialMediaUser{
		ID:            primitive.NewObjectID(),
		Link:          input.Link,
		SocialMediaID: input.SocialMediaID,
		UserID:        input.UserID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		IsActive:      true,
	}

	result, err := collectionSocialMediaUser.InsertOne(ctx, newSocialMediaUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Social media user created",
		"data":    result,
	})
}

// Get Social Media by Auth
func GetMySocialMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get user ID from context
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 2. Validate user existence
	userCollection := database.DB.Collection("users")

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID, "is_active": true}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	// 3. Query social media user
	collection := database.DB.Collection("social_media_user")

	filter := bson.M{
		"user_id": userID,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch social media users",
		})
		return
	}
	defer cursor.Close(ctx)

	var socialMediaUsers []models.SocialMediaUser

	if err := cursor.All(ctx, &socialMediaUsers); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse social media users",
		})
		return
	}

	// 4. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media users fetched successfully",
		"data":    socialMediaUsers,
	})
}

// Get Social Media by User ID
func GetSocialMediaByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	collection := database.DB.Collection("social_media_user")

	filter := bson.M{
		"user_id":   id,
		"is_active": true,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch social media users",
		})
		return
	}
	defer cursor.Close(ctx)

	var socialMediaUsers []models.SocialMediaUser

	if err := cursor.All(ctx, &socialMediaUsers); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse social media users",
		})
		return
	}

	// 4. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media users fetched successfully",
		"data":    socialMediaUsers,
	})
}

func UpdateSocialMediaUser(c *gin.Context) {}

func DeleteSocialMediaUser(c *gin.Context) {}
