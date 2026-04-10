package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort   string
	MongoURI  string
	MongoDB   string
	JWTSecret []byte
}

var Config AppConfig

// Collections name constants
const (
	CollectionUsers           = "users"
	CollectionSocialMedia     = "social_media"
	CollectionSocialMediaUser = "social_media_user"
	CollectionPortofolio      = "portofolio"
	CollectionProject         = "project"
	CollectionProjectImage    = "project_image"
)

func LoadConfig() {
	// Load .env only for local development
	_ = godotenv.Load()

	Config = AppConfig{
		AppPort:   getEnv("APP_PORT", "8081"),
		MongoURI:  getEnv("MONGO_URI", ""),
		MongoDB:   getEnv("MONGO_DB", ""),
		JWTSecret: []byte(getEnv("JWT_SECRET", "default_secret")),
	}

	if Config.MongoURI == "" || Config.MongoDB == "" {
		log.Fatal("❌ MONGO_URI and MONGO_DB must be set")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
