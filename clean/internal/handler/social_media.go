package handler

import (
	"context"
	"net/http"
	"portofolio-api/clean/internal/domain"
	"portofolio-api/clean/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type SocialMediaHandler struct {
	socialMediaService service.SocialMediaService
}

func NewSocialMediaHandler(socialMediaService service.SocialMediaService) *SocialMediaHandler {
	return &SocialMediaHandler{socialMediaService: socialMediaService}
}

func (h *SocialMediaHandler) CreateSocialMedia(c *gin.Context) {
	var input domain.SocialMediaCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	socialMedia, err := h.socialMediaService.Create(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"social_media": socialMedia})
}
