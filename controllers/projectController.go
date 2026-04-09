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

// CreateProjectPortofolio godoc
// @Summary Create a new project for a portfolio
// @Description Create a new project entry associated with a specific portfolio ID owned by the user
// @Tags project
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param project body models.ProjectInput true "Project details"
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project [post]
func CreateProjectPortofolio(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get current user
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	// 2. Bind input
	var input models.ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Check if portofolio belongs to the user
	collectionPortofolio := database.DB.Collection("portofolio")

	portofolioID, err := primitive.ObjectIDFromHex(input.PortofolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Portofolio ID format"})
		return
	}

	err = collectionPortofolio.FindOne(ctx, bson.M{
		"_id":        portofolioID,
		"user_id":    userID,
		"is_deleted": false,
	}).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portofolio not found or not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking portofolio"})
		return
	}

	// 4. Create project
	newProject := models.Project{
		ID:           primitive.NewObjectID(),
		Name:         input.Name,
		Description:  input.Description,
		Link:         input.Link,
		PortofolioID: portofolioID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		IsDeleted:    false,
	}

	collectionProject := database.DB.Collection("project")
	result, err := collectionProject.InsertOne(ctx, newProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"data":    result,
	})
}

// GetProjectsByPortofolioID godoc
// @Summary Get all projects by Portfolio ID
// @Description Retrieve a list of all projects associated with a specific portfolio ID
// @Tags project
// @Produce  json
// @Param portofolio_id path string true "Portfolio ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/portofolio/{portofolio_id} [get]
func GetProjectsByPortofolioID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get portofolio_id from param
	portofolioIDParam := c.Param("portofolio_id")
	portofolioID, err := primitive.ObjectIDFromHex(portofolioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	// 2. Query project
	collection := database.DB.Collection("project")

	filter := bson.M{
		"portofolio_id": portofolioID,
		"is_deleted":    false,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse projects",
		})
		return
	}

	// 3. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Projects fetched successfully",
		"data":    projects,
	})
}

// UpdateProject godoc
// @Summary Update project details
// @Description Update project name, description, and link by project ID for the owner of the associated portfolio
// @Tags project
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param project body models.ProjectUpdateInput true "Updated project details"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/{project_id} [put]
func UpdateProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get current user
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	// 2. Get project_id from param
	projectIDParam := c.Param("project_id")
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 3. Bind input
	var input models.ProjectUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")

	// 4. Find project and check ownership
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": projectID, "is_deleted": false}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project"})
		return
	}

	// Check if portfolio belongs to user
	err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying ownership"})
		return
	}

	// 5. Update project
	update := bson.M{
		"$set": bson.M{
			"name":        input.Name,
			"description": input.Description,
			"link":        input.Link,
			"updated_at":  time.Now(),
		},
	}

	_, err = projectCollection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// UpdateProjectStatus godoc
// @Summary Update project active status
// @Description Toggle the active status of a project by project ID for the owner of the associated portfolio
// @Tags project
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param status body models.ProjectStatusInput true "Updated project status"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/{project_id}/status [patch]
func UpdateProjectStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get current user
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	// 2. Get project_id from param
	projectIDParam := c.Param("project_id")
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 3. Bind input
	var input models.ProjectStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")

	// 4. Find project and check ownership
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": projectID, "is_deleted": false}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project"})
		return
	}

	// Check if portfolio belongs to user
	err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this project status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying ownership"})
		return
	}

	// 5. Update status
	update := bson.M{
		"$set": bson.M{
			"is_active":  input.IsActive,
			"updated_at": time.Now(),
		},
	}

	_, err = projectCollection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project status updated successfully"})
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Soft delete a project entry by project ID for the owner of the associated portfolio
// @Tags project
// @Produce  json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/{project_id} [delete]
func DeleteProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get current user
	user, statusCode, err := helpers.GetCurrentUser(c, ctx)
	if err != nil {
		if statusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized / User not found"})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userID := user.ID

	// 2. Get project_id from param
	projectIDParam := c.Param("project_id")
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")

	// 3. Find project and check ownership
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": projectID, "is_deleted": false}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project"})
		return
	}

	// Check if portfolio belongs to user
	err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying ownership"})
		return
	}

	// 4. Soft Delete
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err = projectCollection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
