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

	"github.com/gin-contrib/cors"
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

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://api-portofolio.declarationdigital.tech"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger route
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(r)
	routes.SocialMediaRoutes(r)
	routes.SocialMediaUserRoutes(r)
	routes.PortofolioRoutes(r)
	routes.ProjectRoutes(r)
	routes.EducationRoutes(r)
	routes.CertificateRoutes(r)
	routes.SkillRoutes(r)
	routes.JobHistoryRoutes(r)

	// Serve static files
	r.Static("/public", "./public")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(":" + port)
}
