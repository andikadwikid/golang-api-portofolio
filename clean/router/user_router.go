package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portofolio-api/clean/internal/handler"
	"portofolio-api/clean/internal/middleware"
)

func UserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	// Health check global (biasanya di root)
	r.GET("/health", healthCheck)

	// Public routes
	public := r.Group("/users")
	{
		public.GET("/", userHandler.GetAllUsers)
		public.GET("/test", testHandler)

		public.POST("/register", userHandler.RegisterUser)
		public.POST("/login", userHandler.LoginUser)
	}

	// Protected routes
	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/:id", userHandler.UpdateUser)
		protected.DELETE("/:id", userHandler.DeleteUser)
	}

}

// handler kecil dipisah (biar reusable & clean)
func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
