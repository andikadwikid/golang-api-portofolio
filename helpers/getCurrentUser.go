package helpers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/database"
	"portofolio-api/models"
)

func GetCurrentUser(c *gin.Context, ctx context.Context) (*models.User, int, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	var user models.User
	err = database.DB.Collection("users").
		FindOne(ctx, bson.M{"_id": userID}).
		Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http.StatusUnauthorized, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return &user, http.StatusOK, nil
}
