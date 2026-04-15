package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"portofolio-api/clean/internal/domain"
	"portofolio-api/clean/internal/pkg/utils"
	"portofolio-api/clean/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input domain.CreateUserInput

	// Bind & validate — tetap di sini, ini memang urusan HTTP layer
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	// Serahkan semua logika ke Service — SATU baris menggantikan ~40 baris
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userService.Register(ctx, input)
	if err != nil {
		// Bedakan error bisnis vs error server
		if err.Error() == "email already registered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      user.ID.Hex(),
	})
}
