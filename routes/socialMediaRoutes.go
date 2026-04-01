package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func SocialMediaRoutes(r *gin.Engine) {
	socialMedia := r.Group("/social-media")
	{
		socialMedia.GET("/", middlewares.AuthMiddleware(), controllers.GetSocialMedia)
		socialMedia.POST("/", middlewares.AuthMiddleware(), controllers.CreateSocialMedia)
		// socialMedia.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateSocialMedia)
		// socialMedia.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteSocialMedia)
	}
}
