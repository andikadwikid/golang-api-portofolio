package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func UserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
		users.GET("/", middlewares.AuthMiddleware(), controllers.GetUsers)
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.LoginUser)
		users.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
		users.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)
	}
}
