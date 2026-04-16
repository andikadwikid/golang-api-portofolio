package router

import (
	"portofolio-api/clean/internal/handler"

	"github.com/gin-gonic/gin"
)

func SocialMediaRouter(r *gin.Engine, socialMediaHandler *handler.SocialMediaHandler) {
	r.POST("/social-media", socialMediaHandler.CreateSocialMedia)
}
