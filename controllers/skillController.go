package controllers

import (
	"context"
	"net/http"
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

// CreateSkill godoc
// @Summary Create a new skill entry
// @Description Create a new skill entry for a user
// @Tags skill
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param skill body models.CreateSkillInput true "Skill creation details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /skill [post]
func CreateSkill(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	var input models.CreateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	collection := database.DB.Collection("skill")

	skill := models.Skill{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, skill)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Skill created successfully",
		"id":      insertedID.Hex(),
	})
}

// GetSkillByID godoc
// @Summary Get skill entry by ID
// @Description Retrieve a single skill entry by its ID
// @Tags skill
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Success 200 {object} models.Skill
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /skill/{id} [get]
func GetSkillByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("skill")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var skill models.Skill
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&skill)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skill)
}

// GetSkillsByUserID godoc
// @Summary Get all skill entries for a user
// @Description Retrieve a list of all skill entries for a specific user ID
// @Tags skill
// @Produce  json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Skill
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /skill/user/{user_id} [get]
func GetSkillsByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("skill")

	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching skill entries"})
		return
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err := cursor.All(ctx, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing skill entries"})
		return
	}

	c.JSON(http.StatusOK, skills)
}

// UpdateSkill godoc
// @Summary Update a skill entry
// @Description Update a skill entry by ID
// @Tags skill
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Param skill body models.UpdateSkillInput true "Updated skill details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /skill/{id} [put]
func UpdateSkill(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.UpdateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	setFields := bson.M{
		"updated_at": time.Now(),
	}

	if input.Name != nil {
		setFields["name"] = *input.Name
	}

	if len(setFields) == 1 { // Only updated_at is present
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	update := bson.M{"$set": setFields}

	result, err := database.DB.Collection("skill").UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill entry not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill entry updated"})
}

// DeleteSkill godoc
// @Summary Delete a skill entry
// @Description Delete a skill entry by ID
// @Tags skill
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /skill/{id} [delete]
func DeleteSkill(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := database.DB.Collection("skill").DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill entry not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill entry deleted"})
}
