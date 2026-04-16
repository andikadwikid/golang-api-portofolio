package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"portofolio-api/clean/config"
	"portofolio-api/clean/internal/handler"
	"portofolio-api/clean/internal/repository"
	"portofolio-api/clean/internal/service"
	"portofolio-api/clean/router"
)

func main() {
	godotenv.Load()

	config.Connect()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	socialMediaRepo := repository.NewSocialMediaRepository()
	socialMediaService := service.NewSocialMediaService(socialMediaRepo)
	socialMediaHandler := handler.NewSocialMediaHandler(socialMediaService)

	r := gin.Default()
	router.UserRoutes(r, userHandler)
	router.SocialMediaRouter(r, socialMediaHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	log.Fatal(r.Run(":" + port))
}
