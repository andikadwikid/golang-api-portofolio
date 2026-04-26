package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func ProjectRoutes(r *gin.Engine) {
	project := r.Group("/project")
	{
		// Project
		project.POST("", middlewares.AuthMiddleware(), controllers.CreateProjectPortofolio)
		project.POST("/", middlewares.AuthMiddleware(), controllers.CreateProjectPortofolio)
		project.GET("/portofolio/:portofolio_id", controllers.GetProjectsByPortofolioID)
		project.PUT("/:project_id", middlewares.AuthMiddleware(), controllers.UpdateProject)
		project.PATCH("/:project_id/status", middlewares.AuthMiddleware(), controllers.UpdateProjectStatus)
		project.DELETE("/:project_id", middlewares.AuthMiddleware(), controllers.DeleteProject)

		// Project Images
		project.GET("/images/my", middlewares.AuthMiddleware(), controllers.GetMyProjectImages)
		project.POST("/:project_id/images", middlewares.AuthMiddleware(), controllers.CreateProjectImage)
		project.GET("/:project_id/images", controllers.GetProjectImagesByProjectID)
		project.PUT("/images/:image_id", middlewares.AuthMiddleware(), controllers.UpdateProjectImage)
		project.PATCH("/images/:image_id/status", middlewares.AuthMiddleware(), controllers.UpdateProjectImageStatus)
		project.DELETE("/images/:image_id", middlewares.AuthMiddleware(), controllers.DeleteProjectImage)
	}
}
