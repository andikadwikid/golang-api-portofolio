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

func (h *UserHandler) LoginUser(c *gin.Context) {
	var input domain.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := h.userService.Login(ctx, input)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input domain.UpdateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.userService.Update(ctx, id, input)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid id" || err.Error() == "no fields to update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.userService.Delete(ctx, id)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid id" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
