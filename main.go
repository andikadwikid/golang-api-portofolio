package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"portofolio-api/config"
	"portofolio-api/database"
	_ "portofolio-api/docs"
	"portofolio-api/routes"
)

func main() {
	// 1. Load configuration
	config.LoadConfig()

	// 2. Connect to Database
	database.Connect()

	// 3. Set up Router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Swagger route
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(r)
	routes.SocialMediaRoutes(r)
	routes.SocialMediaUserRoutes(r)
	routes.PortofolioRoutes(r)
	routes.ProjectRoutes(r)

	// Serve static files
	r.Static("/public", "./public")

	// 4. Graceful Shutdown Implementation
	port := config.Config.AppPort
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("🚀 Server is running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
