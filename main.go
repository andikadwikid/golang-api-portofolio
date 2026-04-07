package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"portofolio-api/database"
	"portofolio-api/routes"
)

func main() {
	// Memuat .env jika ada (biasanya untuk pengembangan lokal)
	// Di server/Docker, variabel lingkungan biasanya sudah diatur via docker-compose atau env_file
	godotenv.Load()

	database.Connect()

	r := gin.Default()
	routes.UserRoutes(r)
	routes.SocialMediaRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(":" + port)
}
