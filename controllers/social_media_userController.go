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

// CreateSocialMediaUser godoc
// @Summary Create a new social media user link
// @Description Link an authenticated user to a specific social media platform with a link
// @Tags social_media_user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param social_media_user body models.SocialMediaUserInput true "Social media user link details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user [post]
func CreateSocialMediaUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized / User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	userID := user.ID

	var input models.SocialMediaUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collectionSocialMedia := database.DB.Collection("social_media")

	socialMediaID, err := primitive.ObjectIDFromHex(input.SocialMediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Social Media ID format"})
		return
	}

	err = collectionSocialMedia.FindOne(ctx, bson.M{"_id": socialMediaID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Social media not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking social media"})
		return
	}

	newSocialMediaUser := models.SocialMediaUser{
		ID:            primitive.NewObjectID(),
		Link:          input.Link,
		SocialMediaID: socialMediaID,
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		IsActive:      true,
	}

	collectionSocialMediaUser := database.DB.Collection("social_media_user")
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

// GetMySocialMedia godoc
// @Summary Get authenticated user's social media links
// @Description Retrieve all social media links associated with the currently authenticated user
// @Tags social_media_user
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user/me [get]
func GetMySocialMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get user from context and validate existence
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized / User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	userID := user.ID

	// 2. Query social media user
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

	// 3. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media users fetched successfully",
		"data":    socialMediaUsers,
	})
}

// GetSocialMediaByUserID godoc
// @Summary Get social media links by User ID
// @Description Retrieve all active social media links associated with a specific user ID
// @Tags social_media_user
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user/user/{user_id} [get]
func GetSocialMediaByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := c.Param("user_id")
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

// UpdateSocialMediaUser godoc
// @Summary Update a social media user link
// @Description Update the link or platform for an existing social media user link by ID
// @Tags social_media_user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param socialmedia_id path string true "Social Media User ID"
// @Param social_media_user body models.SocialMediaUserUpdateInput true "Updated link details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user/{socialmedia_id} [put]
func UpdateSocialMediaUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized / User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	// 1. Dapatkan user ID dari konteks autentikasi
	userID := user.ID

	// 2. Ambil ID social media user dari parameter URL
	idParam := c.Param("socialmedia_id")
	socialMediaUserID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 3. Bind dan validasi input JSON
	var input models.SocialMediaUserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.DB.Collection("social_media_user")

	// 4. Validasi dan pastikan social media yang direferensikan ada di database
	socialMediaID, err := primitive.ObjectIDFromHex(input.SocialMediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Social Media ID format"})
		return
	}

	collectionSocialMedia := database.DB.Collection("social_media")
	err = collectionSocialMedia.FindOne(ctx, bson.M{"_id": socialMediaID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Social media not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking social media"})
		return
	}

	// 5. Siapkan data yang akan diperbarui
	update := bson.M{
		"$set": bson.M{
			"link":            input.Link,
			"social_media_id": socialMediaID,
			"updated_at":      time.Now(),
		},
	}

	// 6. Jalankan update dengan filter ID dan UserID (ownership check)
	result, err := collection.UpdateOne(ctx, bson.M{"_id": socialMediaUserID, "user_id": userID}, update)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media user"})
		return
	}

	// 7. Periksa apakah ada dokumen yang cocok dan berhasil diupdate
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media user not found or not authorized"})
		return
	}

	// 8. Berikan respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media user updated successfully",
	})
}

// UpdateSocialMediaUserStatus godoc
// @Summary Toggle social media user link active status
// @Description Update the active status of a social media user link by ID for the owner
// @Tags social_media_user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param socialmedia_id path string true "Social Media User ID"
// @Param status body models.SocialMediaUserStatusUpdateInput true "Updated status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user/{socialmedia_id}/status [patch]
func UpdateSocialMediaUserStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized / User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	userID := user.ID

	idParam := c.Param("socialmedia_id")
	socialMediaUserID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input models.SocialMediaUserStatusUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.DB.Collection("social_media_user")

	update := bson.M{
		"$set": bson.M{
			"is_active":  input.IsActive,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": socialMediaUserID, "user_id": userID}, update)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media user status"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media user not found or not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Social media user status updated successfully",
	})
}

// DeleteSocialMediaUser godoc
// @Summary Delete a social media user link
// @Description Delete a social media user link entry by ID for the owner
// @Tags social_media_user
// @Produce  json
// @Security BearerAuth
// @Param socialmedia_id path string true "Social Media User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social-media-user/{socialmedia_id} [delete]
func DeleteSocialMediaUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Dapatkan user dari context autentikasi
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized / User not found",
			})
			return
		}

		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	userID := user.ID

	// 2. Ambil ID social media user dari parameter URL (:socialmedia_id)
	idParam := c.Param("socialmedia_id")
	socialMediaUserID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 3. Eksekusi penghapusan dengan filter ID dan UserID (ownership check)
	collection := database.DB.Collection("social_media_user")

	result, err := collection.DeleteOne(ctx, bson.M{
		"_id":     socialMediaUserID,
		"user_id": userID,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media user"})
		return
	}

	// 4. Periksa apakah ada dokumen yang terhapus
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media user not found or not authorized"})
		return
	}

	// 5. Berikan respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Social media user deleted successfully",
	})
}
