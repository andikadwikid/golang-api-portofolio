package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return primitive.NilObjectID, errors.New("unauthorized")
	}

	return primitive.ObjectIDFromHex(userIDStr.(string))
}
