package routes

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/controllers"
	"portofolio-api/middlewares"
)

func SocialMediaUserRoutes(r *gin.Engine) {
	socialMediaUser := r.Group("/social-media-user")
	{
		socialMediaUser.POST("", middlewares.AuthMiddleware(), controllers.CreateSocialMediaUser)
		socialMediaUser.POST("/", middlewares.AuthMiddleware(), controllers.CreateSocialMediaUser)
		socialMediaUser.GET("/me", middlewares.AuthMiddleware(), controllers.GetMySocialMedia)
		socialMediaUser.GET("/user/:user_id", middlewares.AuthMiddleware(), controllers.GetSocialMediaByUserID)
		socialMediaUser.PUT("/:socialmedia_id", middlewares.AuthMiddleware(), controllers.UpdateSocialMediaUser)
		socialMediaUser.PATCH("/:socialmedia_id/status", middlewares.AuthMiddleware(), controllers.UpdateSocialMediaUserStatus)
		socialMediaUser.DELETE("/:socialmedia_id", middlewares.AuthMiddleware(), controllers.DeleteSocialMediaUser)
	}
}
