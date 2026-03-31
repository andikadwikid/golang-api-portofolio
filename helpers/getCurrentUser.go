package helpers

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"portofolio-api/database"
	"portofolio-api/models"
)

func GetCurrentUser(c *gin.Context, ctx context.Context) (*models.User, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = database.DB.Collection("users").
		FindOne(ctx, bson.M{"_id": userID}).
		Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user_id not found in context")
	}

	return userID.(string), nil
}
