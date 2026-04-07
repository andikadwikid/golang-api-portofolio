package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"portofolio-api/database"
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
