package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/database"
	"portofolio-api/helpers"
	"portofolio-api/models"
)

// CreateProjectImage godoc
// @Summary Upload multiple images for a project
// @Description Upload multiple images (jpg, jpeg, png) for a specific project owned by the user
// @Tags project_image
// @Accept  multipart/form-data
// @Produce  json
// @Security BearerAuth
// @Param project_id formData string true "Project ID"
// @Param images formData file true "Project images (multiple)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/{project_id}/images [post]
func CreateProjectImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// 2. Get project_id from form
	projectIDParam := c.Param("project_id") // We'll take it from URL for consistency if route is /project/:project_id/images
	if projectIDParam == "" {
		projectIDParam = c.PostForm("project_id")
	}

	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID format"})
		return
	}

	// 3. Check project ownership
	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")

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

	// Verify portfolio ownership
	err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to add images to this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying ownership"})
		return
	}

	// 4. Handle multiple upload
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
		return
	}

	// Create directory if not exists
	uploadDir := "./public/uploads/project_images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	var newImages []interface{}
	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s has invalid extension. Only jpg, jpeg, png allowed.", file.Filename)})
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), primitive.NewObjectID().Hex(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", file.Filename)})
			return
		}

		// Prepare document
		imageDoc := models.ProjectImage{
			ID:        primitive.NewObjectID(),
			ProjectID: projectID,
			Image:     "public/uploads/project_images/" + filename,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive:  true,
			IsDeleted: false,
		}
		newImages = append(newImages, imageDoc)
	}

	// 5. Bulk insert
	imageCollection := database.DB.Collection("project_image")
	_, err = imageCollection.InsertMany(ctx, newImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image records to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%d images uploaded successfully", len(files)),
		"data":    newImages,
	})
}

// UpdateProjectImage godoc
// @Summary Update a project image
// @Description Update project image file or associated project ID. If a new image is uploaded, the old one is deleted.
// @Tags project_image
// @Accept  multipart/form-data
// @Produce  json
// @Security BearerAuth
// @Param image_id path string true "Image ID"
// @Param image formData file false "New project image"
// @Param project_id formData string false "New Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/images/{image_id} [put]
func UpdateProjectImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Get current user
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

	// 2. Get image_id from param
	imageIDParam := c.Param("image_id")
	imageID, err := primitive.ObjectIDFromHex(imageIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Image ID format"})
		return
	}

	// 3. Find existing image and verify ownership
	imageCollection := database.DB.Collection("project_image")
	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")

	var existingImage models.ProjectImage
	err = imageCollection.FindOne(ctx, bson.M{"_id": imageID, "is_deleted": false}).Decode(&existingImage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching image"})
		return
	}

	// Check current project ownership
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": existingImage.ProjectID, "is_deleted": false}).Decode(&project)
	if err == nil {
		err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to update this image"})
		return
	}

	// 4. Handle updates
	update := bson.M{"updated_at": time.Now()}
	newProjectIDParam := c.PostForm("project_id")
	if newProjectIDParam != "" {
		newProjectID, err := primitive.ObjectIDFromHex(newProjectIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID format"})
			return
		}

		var newProject models.Project
		err = projectCollection.FindOne(ctx, bson.M{"_id": newProjectID, "is_deleted": false}).Decode(&newProject)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project"})
			return
		}

		err = portofolioCollection.FindOne(ctx, bson.M{"_id": newProject.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to associate image with new project"})
			return
		}

		update["project_id"] = newProjectID
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		// Validate extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid extension"})
			return
		}

		// Save new file
		uploadDir := "./public/uploads/project_images"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), primitive.NewObjectID().Hex(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new file"})
			return
		}

		// Delete old file
		if existingImage.Image != "" {
			oldPath := "./" + existingImage.Image
			os.Remove(oldPath)
		}

		update["image"] = "public/uploads/project_images/" + filename
	}

	_, err = imageCollection.UpdateOne(ctx, bson.M{"_id": imageID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project image updated successfully"})
}

// UpdateProjectImageStatus godoc
// @Summary Toggle project image active status
// @Description Update the active status of a project image by ID for the owner of the associated portfolio
// @Tags project_image
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param image_id path string true "Image ID"
// @Param status body models.ProjectImageStatusInput true "Updated status"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /project/images/{image_id}/status [patch]
func UpdateProjectImageStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Get current user
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

	// 2. Get image_id from param
	imageIDParam := c.Param("image_id")
	imageID, err := primitive.ObjectIDFromHex(imageIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Image ID format"})
		return
	}

	// 3. Ownership check
	imageCollection := database.DB.Collection("project_image")
	var existingImage models.ProjectImage
	err = imageCollection.FindOne(ctx, bson.M{"_id": imageID, "is_deleted": false}).Decode(&existingImage)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": existingImage.ProjectID, "is_deleted": false}).Decode(&project)
	if err == nil {
		err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 4. Update status
	var input models.ProjectImageStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = imageCollection.UpdateOne(ctx, bson.M{"_id": imageID}, bson.M{"$set": bson.M{
		"is_active":  input.IsActive,
		"updated_at": time.Now(),
	}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project image status updated successfully"})
}

// DeleteProjectImage godoc
// @Summary Delete a project image
// @Description Soft delete a project image by ID for the owner of the associated portfolio
// @Tags project_image
// @Produce  json
// @Security BearerAuth
// @Param image_id path string true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/images/{image_id} [delete]
func DeleteProjectImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Get current user
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

	// 2. Get image_id from param
	imageIDParam := c.Param("image_id")
	imageID, err := primitive.ObjectIDFromHex(imageIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Image ID format"})
		return
	}

	// 3. Ownership check
	imageCollection := database.DB.Collection("project_image")
	var existingImage models.ProjectImage
	err = imageCollection.FindOne(ctx, bson.M{"_id": imageID, "is_deleted": false}).Decode(&existingImage)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	projectCollection := database.DB.Collection("project")
	portofolioCollection := database.DB.Collection("portofolio")
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{"_id": existingImage.ProjectID, "is_deleted": false}).Decode(&project)
	if err == nil {
		err = portofolioCollection.FindOne(ctx, bson.M{"_id": project.PortofolioID, "user_id": userID, "is_deleted": false}).Err()
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

	_, err = imageCollection.UpdateOne(ctx, bson.M{"_id": imageID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project image deleted successfully"})
}

// GetMyProjectImages godoc
// @Summary Get all my project images
// @Description Retrieve all project images belonging to the authenticated user using aggregation
// @Tags project_image
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/images/my [get]
func GetMyProjectImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Get current user
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

	// 2. Aggregate images
	imageCollection := database.DB.Collection("project_image")

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "project"},
			{Key: "localField", Value: "project_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "project"},
		}}},
		{{Key: "$unwind", Value: "$project"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "portofolio"},
			{Key: "localField", Value: "project.portofolio_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "portofolio"},
		}}},
		{{Key: "$unwind", Value: "$portofolio"}},
		{{Key: "$match", Value: bson.D{
			{Key: "portofolio.user_id", Value: userID},
			{Key: "is_deleted", Value: false},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "project", Value: 0},
			{Key: "portofolio", Value: 0},
		}}},
	}

	cursor, err := imageCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate images"})
		return
	}
	defer cursor.Close(ctx)

	var images []models.ProjectImage
	if err := cursor.All(ctx, &images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User project images fetched successfully",
		"data":    images,
	})
}

// GetProjectImagesByProjectID godoc
// @Summary Get images by Project ID
// @Description Retrieve all active images associated with a specific project ID
// @Tags project_image
// @Produce  json
// @Param project_id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /project/{project_id}/images [get]
func GetProjectImagesByProjectID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectIDParam := c.Param("project_id")
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project ID format"})
		return
	}

	imageCollection := database.DB.Collection("project_image")

	filter := bson.M{
		"project_id": projectID,
		"is_deleted": false,
		"is_active":  true,
	}

	cursor, err := imageCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}
	defer cursor.Close(ctx)

	var images []models.ProjectImage
	if err := cursor.All(ctx, &images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project images fetched successfully",
		"data":    images,
	})
}
