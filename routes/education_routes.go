package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func EducationRoutes(r *gin.Engine) {
	education := r.Group("/education")
	{
		education.POST("/", middlewares.AuthMiddleware(), controllers.CreateEducation)
		education.GET("/:id", middlewares.AuthMiddleware(), controllers.GetEducationByID)
		education.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetEducationByUserID)
		education.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateEducation)
		education.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteEducation)
	}
}
