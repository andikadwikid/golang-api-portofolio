package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"portofolio-api/database"
	"portofolio-api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env")
	}

	database.Connect()

	r := gin.Default()
	routes.UserRoutes(r)
	routes.SocialMediaRoutes(r)

	r.Run(":8081")
}
