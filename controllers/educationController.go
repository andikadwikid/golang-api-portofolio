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
	"portofolio-api/models"
	"portofolio-api/utils"
)

// CreateEducation godoc
// @Summary Create a new education entry
// @Description Create a new education entry for a user
// @Tags education
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param education body models.CreateEducationInput true "Education creation details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /education [post]
func CreateEducation(c *gin.Context) {
	var input models.CreateEducationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	collection := database.DB.Collection("education")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	education := models.Education{
		ID:          primitive.NewObjectID(),
		Name:        input.Name,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Score:       input.Score,
		Description: input.Description,
		Faculty:     input.Faculty,
		Degree:      input.Degree,
		UserID:      input.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, education)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Education created successfully",
		"id":      insertedID.Hex(),
	})
}

// GetEducationByID godoc
// @Summary Get education entry by ID
// @Description Retrieve a single education entry by its ID
// @Tags education
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Education ID"
// @Success 200 {object} models.Education
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /education/{id} [get]
func GetEducationByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("education")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var education models.Education
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&education)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Education not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, education)
}

// GetEducationByUserID godoc
// @Summary Get all education entries for a user
// @Description Retrieve a list of all education entries for a specific user ID
// @Tags education
// @Produce  json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Education
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /education/user/{user_id} [get]
func GetEducationByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("education")

	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching education entries"})
		return
	}
	defer cursor.Close(ctx)

	var education []models.Education
	if err := cursor.All(ctx, &education); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing education entries"})
		return
	}

	c.JSON(http.StatusOK, education)
}

// UpdateEducation godoc
// @Summary Update an education entry
// @Description Update an education entry by ID
// @Tags education
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Education ID"
// @Param education body models.UpdateEducationInput true "Updated education details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /education/{id} [put]
func UpdateEducation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("education")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.UpdateEducationInput
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
	if input.StartDate != nil {
		setFields["start_date"] = *input.StartDate
	}
	if input.EndDate != nil {
		setFields["end_date"] = *input.EndDate
	}
	if input.Score != nil {
		setFields["score"] = *input.Score
	}
	if input.Description != nil {
		setFields["description"] = *input.Description
	}
	if input.Faculty != nil {
		setFields["faculty"] = *input.Faculty
	}
	if input.Degree != nil {
		setFields["degree"] = *input.Degree
	}

	if len(setFields) == 1 { // Only updated_at is present
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	update := bson.M{"$set": setFields}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Education entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education entry updated"})
}

// DeleteEducation godoc
// @Summary Delete an education entry
// @Description Delete an education entry by ID
// @Tags education
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Education ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /education/{id} [delete]
func DeleteEducation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("education")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Education entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education entry deleted"})
}
