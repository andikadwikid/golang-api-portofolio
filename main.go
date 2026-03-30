package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"portofolio-api/database"
	"portofolio-api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	r := gin.Default()
	routes.UserRoutes(r)

	r.Run(":8080")
}
