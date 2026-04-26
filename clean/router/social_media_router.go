package router

import (
	"github.com/gin-gonic/gin"

	"portofolio-api/clean/internal/handler"
)

func SocialMediaRouter(r *gin.Engine, socialMediaHandler *handler.SocialMediaHandler) {
	r.POST("/social-media", socialMediaHandler.CreateSocialMedia)
	r.GET("/social-media", socialMediaHandler.GetAllSocialMedia)
	r.GET("/social-media/:id", socialMediaHandler.GetSocialMediaById)
	r.PUT("/social-media/:id", socialMediaHandler.UpdateSocialMedia)
	r.DELETE("/social-media/:id", socialMediaHandler.DeleteSocialMedia)
}
