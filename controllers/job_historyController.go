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

// CreateJobHistory godoc
// @Summary Create a new job history entry
// @Description Create a new job history entry for a user
// @Tags job_history
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param job_history body models.CreateJobHistoryInput true "Job History creation details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-history [post]
func CreateJobHistory(c *gin.Context) {
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

	var input models.CreateJobHistoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	collection := database.DB.Collection("job_history")

	var skillObjectIDs []primitive.ObjectID
	for _, skillIDStr := range input.SkillIDs {
		skillID, err := primitive.ObjectIDFromHex(skillIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Skill ID format"})
			return
		}
		skillObjectIDs = append(skillObjectIDs, skillID)
	}

	jobHistory := models.JobHistory{
		ID:            primitive.NewObjectID(),
		CompanyName:   input.CompanyName,
		Position:      input.Position,
		JobDescription: input.JobDescription,
		Duration:      input.Duration,
		UserID:        userID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		IsCurrent:     input.IsCurrent,
		SkillIDs:      skillObjectIDs,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	result, err := collection.InsertOne(ctx, jobHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Job History created successfully",
		"id":      insertedID.Hex(),
	})
}

// GetJobHistoryByID godoc
// @Summary Get job history entry by ID
// @Description Retrieve a single job history entry by its ID
// @Tags job_history
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Job History ID"
// @Success 200 {object} models.JobHistory
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-history/{id} [get]
func GetJobHistoryByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("job_history")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var jobHistory models.JobHistory
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&jobHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job History not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobHistory)
}

// GetJobHistoryByUserID godoc
// @Summary Get all job history entries for a user
// @Description Retrieve a list of all job history entries for a specific user ID
// @Tags job_history
// @Produce  json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} models.JobHistory
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-history/user/{user_id} [get]
func GetJobHistoryByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("job_history")

	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching job history entries"})
		return
	}
	defer cursor.Close(ctx)

	var jobHistories []models.JobHistory
	if err := cursor.All(ctx, &jobHistories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing job history entries"})
		return
	}

	c.JSON(http.StatusOK, jobHistories)
}

// UpdateJobHistory godoc
// @Summary Update a job history entry
// @Description Update a job history entry by ID
// @Tags job_history
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Job History ID"
// @Param job_history body models.UpdateJobHistoryInput true "Updated job history details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-history/{id} [put]
func UpdateJobHistory(c *gin.Context) {
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

	var input models.UpdateJobHistoryInput
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

	if input.CompanyName != nil {
		setFields["company_name"] = *input.CompanyName
	}
	if input.Position != nil {
		setFields["position"] = *input.Position
	}
	if input.JobDescription != nil {
		setFields["job_description"] = *input.JobDescription
	}
	if input.Duration != nil {
		setFields["duration"] = *input.Duration
	}
	if input.StartDate != nil {
		setFields["start_date"] = *input.StartDate
	}
	if input.EndDate != nil {
		setFields["end_date"] = *input.EndDate
	}
	if input.IsCurrent != nil {
		setFields["is_current"] = *input.IsCurrent
	}
	if input.SkillIDs != nil {
		var skillObjectIDs []primitive.ObjectID
		for _, skillIDStr := range input.SkillIDs {
			skillID, err := primitive.ObjectIDFromHex(skillIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Skill ID format"})
				return
			}
			skillObjectIDs = append(skillObjectIDs, skillID)
		}
		setFields["skill_ids"] = skillObjectIDs
	}

	if len(setFields) == 1 { // Only updated_at is present
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	update := bson.M{"$set": setFields}

	result, err := database.DB.Collection("job_history").UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job History entry not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job History entry updated"})
}

// DeleteJobHistory godoc
// @Summary Delete a job history entry
// @Description Delete a job history entry by ID
// @Tags job_history
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Job History ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /job-history/{id} [delete]
func DeleteJobHistory(c *gin.Context) {
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

	result, err := database.DB.Collection("job_history").DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job History entry not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job History entry deleted"})
}
