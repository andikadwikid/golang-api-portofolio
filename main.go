package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"portofolio-api/database"
	_ "portofolio-api/docs"
	"portofolio-api/routes"
)

// @title Portfolio API
// @version 1.0
// @description This is a portfolio API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Memuat .env jika ada (biasanya untuk pengembangan lokal)
	// Di server/Docker, variabel lingkungan biasanya sudah diatur via docker-compose atau env_file
	godotenv.Load()

	database.Connect()

	r := gin.Default()

	// Swagger route
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(r)
	routes.SocialMediaRoutes(r)
	routes.SocialMediaUserRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(":" + port)
}
