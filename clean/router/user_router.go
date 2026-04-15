package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portofolio-api/clean/internal/handler"
)

func UserRoutes(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	users := r.Group("/users")
	{
		users.GET("/", userHandler.GetAllUsers)
		users.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
		users.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
		// users.GET("/", middlewares.AuthMiddleware(), controllers.GetUsers)
		users.POST("/register", userHandler.RegisterUser)
		users.POST("/login", userHandler.LoginUser)
	}

	return r
}
