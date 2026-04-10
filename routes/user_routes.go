package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/database"
	"portofolio-api/middlewares"
	"portofolio-api/repositories"
)

func UserRoutes(r *gin.Engine) {
	// Initialize Repository and Controller
	userRepo := repositories.NewUserRepository(database.DB)
	userCtrl := controllers.NewUserController(userRepo)

	users := r.Group("/users")
	{
		users.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
		users.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
		users.GET("/", middlewares.AuthMiddleware(), userCtrl.GetUsers)
		users.POST("/register", userCtrl.RegisterUser)
		users.POST("/login", userCtrl.LoginUser)
		users.PUT("/:id", middlewares.AuthMiddleware(), userCtrl.UpdateUser)
		users.DELETE("/:id", middlewares.AuthMiddleware(), userCtrl.DeleteUser)
	}
}
