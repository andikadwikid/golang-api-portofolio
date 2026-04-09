package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"portofolio-api/database"
	"portofolio-api/helpers"
	"portofolio-api/models"
)

// CreatePortofolio godoc
// @Summary Create a new portofolio
// @Description Create a new portofolio entry for the authenticated user
// @Tags portofolio
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param portofolio body models.PortofolioInput true "Portofolio details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio [post]
func CreatePortofolio(c *gin.Context) {
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

	var input models.PortofolioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPortofolio := models.Portofolio{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
		IsDeleted: false,
	}

	collectionPortofolio := database.DB.Collection("portofolio")
	result, err := collectionPortofolio.InsertOne(ctx, newPortofolio)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Portofolio created successfully",
		"data":    result,
	})
}

// GetMyPortofolios godoc
// @Summary Get all portofolios by Auth user
// @Description Get all portofolios for the authenticated user where is_deleted is false
// @Tags portofolio
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio [get]
func GetMyPortofolios(c *gin.Context) {
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

	// 2. Query portofolio
	collectionPortofolio := database.DB.Collection("portofolio")

	filter := bson.M{
		"user_id":    userID,
		"is_deleted": false,
	}

	cursor, err := collectionPortofolio.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch portofolios",
		})
		return
	}
	defer cursor.Close(ctx)

	var portofolios []models.Portofolio

	if err := cursor.All(ctx, &portofolios); err != nil {
		log.Println("Error decoding portofolios:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse portofolios",
		})
		return
	}

	log.Printf("Fetched %d portofolios for user %s: %+v\n", len(portofolios), userID.Hex(), portofolios)

	// 3. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolios fetched successfully",
		"data":    portofolios,
	})
}

// GetPortofolioByUserID godoc
// @Summary Get all portofolios by User ID
// @Description Get all active portofolios for a specific user ID where is_deleted is false
// @Tags portofolio
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio/user/{user_id} [get]
func GetPortofolioByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get user_id from param
	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	// 2. Query portofolio
	collectionPortofolio := database.DB.Collection("portofolio")

	filter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"is_deleted": false,
	}

	cursor, err := collectionPortofolio.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch portofolios",
		})
		return
	}
	defer cursor.Close(ctx)

	var portofolios []models.Portofolio = []models.Portofolio{}

	if err := cursor.All(ctx, &portofolios); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse portofolios",
		})
		return
	}

	// 3. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolios fetched successfully",
		"data":    portofolios,
	})
}

// UpdatePortofolio godoc
// @Summary Update a portofolio entry
// @Description Update portfolio details by ID for the authenticated user
// @Tags portofolio
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param portofolio_id path string true "Portofolio ID"
// @Param portofolio body models.PortofolioInput true "Updated portofolio details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio/{portofolio_id} [put]
func UpdatePortofolio(c *gin.Context) {
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

	// 2. Get portfolio ID from param
	portofolioIDParam := c.Param("portofolio_id")
	portofolioID, err := primitive.ObjectIDFromHex(portofolioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	// 3. Bind input
	var input models.PortofolioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Update portofolio
	collectionPortofolio := database.DB.Collection("portofolio")

	update := bson.M{
		"$set": bson.M{
			"name":       input.Name,
			"updated_at": time.Now(),
		},
	}

	result, err := collectionPortofolio.UpdateOne(ctx, bson.M{"_id": portofolioID, "user_id": userID, "is_deleted": false}, update)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update portofolio",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Portofolio not found or not authorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolio updated successfully",
	})
}

// UpdatePortofolioStatus godoc
// @Summary Update portofolio status (is_active)
// @Description Toggle the active status of a portofolio entry by ID for the authenticated user
// @Tags portofolio
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param portofolio_id path string true "Portofolio ID"
// @Param status body models.PortofolioStatusInput true "Updated portofolio status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio/{portofolio_id}/status [patch]
func UpdatePortofolioStatus(c *gin.Context) {
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

	// 2. Get portfolio ID from param
	portofolioIDParam := c.Param("portofolio_id")
	portofolioID, err := primitive.ObjectIDFromHex(portofolioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	// 3. Bind input
	var input models.PortofolioStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Update status
	collectionPortofolio := database.DB.Collection("portofolio")

	update := bson.M{
		"$set": bson.M{
			"is_active":  input.IsActive,
			"updated_at": time.Now(),
		},
	}

	result, err := collectionPortofolio.UpdateOne(ctx, bson.M{"_id": portofolioID, "user_id": userID, "is_deleted": false}, update)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update portofolio status",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Portofolio not found or not authorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolio status updated successfully",
	})
}

// DeletePortofolio godoc
// @Summary Delete a portofolio entry
// @Description Soft delete a portofolio entry by ID for the authenticated user
// @Tags portofolio
// @Produce  json
// @Security BearerAuth
// @Param portofolio_id path string true "Portofolio ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio/{portofolio_id} [delete]
func DeletePortofolio(c *gin.Context) {
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

	// 2. Get portfolio ID from param
	portofolioIDParam := c.Param("portofolio_id")
	portofolioID, err := primitive.ObjectIDFromHex(portofolioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	// 3. Update status (Soft Delete)
	collectionPortofolio := database.DB.Collection("portofolio")

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	result, err := collectionPortofolio.UpdateOne(ctx, bson.M{"_id": portofolioID, "user_id": userID, "is_deleted": false}, update)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete portofolio",
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Portofolio not found or not authorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolio deleted successfully",
	})
}

// GetAllPortofolios godoc
// @Summary Get all active portfolios
// @Description Retrieve a list of all portfolios that are marked as active and not deleted
// @Tags portofolio
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /portofolio/all [get]
func GetAllPortofolios(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionPortofolio := database.DB.Collection("portofolio")

	// Filter: must be active and not deleted
	filter := bson.M{
		"is_active":  true,
		"is_deleted": false,
	}

	cursor, err := collectionPortofolio.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch portofolios",
		})
		return
	}
	defer cursor.Close(ctx)

	var portofolios []models.Portofolio
	if err := cursor.All(ctx, &portofolios); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse portofolios",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portofolios fetched successfully",
		"data":    portofolios,
	})
}
