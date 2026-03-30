package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func UserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.LoginUser)
		users.GET("/", middlewares.AuthMiddleware(), controllers.GetUsers)
		users.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
		users.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)
	}
}
