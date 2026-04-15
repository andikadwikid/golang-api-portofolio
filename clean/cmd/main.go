package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"portofolio-api/clean/config"
	"portofolio-api/clean/internal/handler"
	"portofolio-api/clean/internal/repository"
	"portofolio-api/clean/internal/service"
	"portofolio-api/clean/router"
)

func main() {
	godotenv.Load()

	// 1. Koneksi DB
	config.Connect()

	// 2. Repository — hanya tahu MongoDB
	userRepo := repository.NewUserRepository(config.DB.Collection("users"))

	// 3. Service — hanya tahu Repository
	userService := service.NewUserService(userRepo)

	// 4. Handler — hanya tahu Service
	userHandler := handler.NewUserHandler(userService)

	r := router.UserRoutes(userHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
