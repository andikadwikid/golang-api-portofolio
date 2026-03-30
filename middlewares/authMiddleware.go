package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"portofolio-api/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user ID to context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
