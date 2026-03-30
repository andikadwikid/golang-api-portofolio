package helpers

import (
	"context"

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
