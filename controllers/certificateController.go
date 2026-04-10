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
	"portofolio-api/utils"
)

// CreateCertificate godoc
// @Summary Create a new certificate entry
// @Description Create a new certificate entry for a user with file upload
// @Tags certificate
// @Accept  multipart/form-data
// @Produce  json
// @Security BearerAuth
// @Param title formData string true "Certificate Title"
// @Param organizer formData string true "Certificate Organizer"
// @Param year formData string true "Certificate Year"
// @Param description formData string false "Certificate Description"
// @Param file formData file true "Certificate file (pdf, png, jpg, jpeg)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /certificate [post]
func CreateCertificate(c *gin.Context) {
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

	// 2. Bind form data
	var input models.CreateCertificateInput
	if err := c.ShouldBind(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	// 3. Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true}
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s has invalid extension. Only pdf, png, jpg, jpeg allowed.", file.Filename)})
		return
	}

	// Create directory if not exists
	uploadDir := "./public/uploads/certificates"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), primitive.NewObjectID().Hex(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", file.Filename)})
		return
	}

	// 4. Save certificate details to database
	collection := database.DB.Collection("certificate")

	certificate := models.Certificate{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Organizer:   input.Organizer,
		UrlFile:     "public/uploads/certificates/" + filename,
		Year:        input.Year,
		Description: input.Description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, certificate)
	if err != nil {
		// If database insertion fails, delete the uploaded file
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Certificate created successfully",
		"id":      insertedID.Hex(),
	})
}

// GetCertificateByID godoc
// @Summary Get certificate entry by ID
// @Description Retrieve a single certificate entry by its ID
// @Tags certificate
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Certificate ID"
// @Success 200 {object} models.Certificate
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /certificate/{id} [get]
func GetCertificateByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("certificate")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var certificate models.Certificate
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&certificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certificate)
}

// GetCertificatesByUserID godoc
// @Summary Get all certificate entries for a user
// @Description Retrieve a list of all certificate entries for a specific user ID
// @Tags certificate
// @Produce  json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Certificate
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /certificate/user/{user_id} [get]
func GetCertificatesByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Collection("certificate")

	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching certificate entries"})
		return
	}
	defer cursor.Close(ctx)

	var certificates []models.Certificate
	if err := cursor.All(ctx, &certificates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing certificate entries"})
		return
	}

	c.JSON(http.StatusOK, certificates)
}

// UpdateCertificate godoc
// @Summary Update a certificate entry
// @Description Update a certificate entry by ID, optionally replacing the file
// @Tags certificate
// @Accept  multipart/form-data
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Certificate ID"
// @Param title formData string false "Certificate Title"
// @Param organizer formData string false "Certificate Organizer"
// @Param year formData string false "Certificate Year"
// @Param description formData string false "Certificate Description"
// @Param file formData file false "New certificate file (pdf, png, jpg, jpeg)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /certificate/{id} [put]
func UpdateCertificate(c *gin.Context) {
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

	// 2. Get certificate_id from param
	certificateIDParam := c.Param("id")
	certificateID, err := primitive.ObjectIDFromHex(certificateIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Certificate ID format"})
		return
	}

	// 3. Find existing certificate and verify ownership
	collection := database.DB.Collection("certificate")
	var existingCertificate models.Certificate
	err = collection.FindOne(ctx, bson.M{"_id": certificateID, "user_id": userID}).Decode(&existingCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found or unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching certificate"})
		return
	}

	// 4. Handle updates
	updateFields := bson.M{"updated_at": time.Now()}

	if title := c.PostForm("title"); title != "" {
		updateFields["title"] = title
	}
	if organizer := c.PostForm("organizer"); organizer != "" {
		updateFields["organizer"] = organizer
	}
	if year := c.PostForm("year"); year != "" {
		updateFields["year"] = year
	}
	if description := c.PostForm("description"); description != "" {
		updateFields["description"] = description
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err == nil { // A new file was uploaded
		// Validate file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExtensions := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true}
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s has invalid extension. Only pdf, png, jpg, jpeg allowed.", file.Filename)})
			return
		}

		// Save new file
		uploadDir := "./public/uploads/certificates"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), primitive.NewObjectID().Hex(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save new file %s", file.Filename)})
			return
		}

		// Delete old file
		if existingCertificate.UrlFile != "" {
			oldPath := "./" + existingCertificate.UrlFile
			os.Remove(oldPath)
		}

		updateFields["url_file"] = "public/uploads/certificates/" + filename
	}

	if len(updateFields) == 1 { // Only updated_at is present
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": certificateID}, bson.M{"$set": updateFields})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate updated successfully"})
}

// DeleteCertificate godoc
// @Summary Delete a certificate entry
// @Description Delete a certificate entry by ID for the owner
// @Tags certificate
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Certificate ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /certificate/{id} [delete]
func DeleteCertificate(c *gin.Context) {
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

	// 2. Get certificate_id from param
	certificateIDParam := c.Param("id")
	certificateID, err := primitive.ObjectIDFromHex(certificateIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Certificate ID format"})
		return
	}

	// 3. Find existing certificate and verify ownership
	collection := database.DB.Collection("certificate")
	var existingCertificate models.Certificate
	err = collection.FindOne(ctx, bson.M{"_id": certificateID, "user_id": userID}).Decode(&existingCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found or unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching certificate"})
		return
	}

	// 4. Delete file from server
	if existingCertificate.UrlFile != "" {
		filePath := "./" + existingCertificate.UrlFile
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to delete file %s: %v", filePath, err)
			// Continue with database deletion even if file deletion fails
		}
	}

	// 5. Delete certificate from database
	result, err := collection.DeleteOne(ctx, bson.M{"_id": certificateID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
