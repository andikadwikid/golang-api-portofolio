package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func SocialMediaUserRoutes(r *gin.Engine) {
	socialMediaUser := r.Group("/social-media-user")
	{
		socialMediaUser.GET("/me", middlewares.AuthMiddleware(), controllers.GetMySocialMedia)
		socialMediaUser.GET("/user/:id", middlewares.AuthMiddleware(), controllers.GetSocialMediaByUserID)
		socialMediaUser.POST("/", middlewares.AuthMiddleware(), controllers.CreateSocialMediaUser)
		socialMediaUser.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateSocialMediaUser)
		socialMediaUser.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteSocialMediaUser)
	}
}
