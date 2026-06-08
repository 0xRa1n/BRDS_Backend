package main

import (
	"log"
	"os"

	"backend/config"
	"backend/middleware"
	"backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load dotenv
	dotenvErr := godotenv.Load(".env")
	if dotenvErr != nil {
		log.Fatalf("Error loading .env file: %s", dotenvErr)
	}


	// Initialize In-Memory Cache
	config.InitCache()

	// Initialize Database Connection
	// Note: You must have PostgreSQL running and credentials configured via environment variables
	// DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT
	config.InitDB()

	// Initialize Gin Router
	router := gin.Default()

	// Enable CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"} // Adjust depending on frontend port
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "Idempotency-Key"}
	router.Use(cors.New(corsConfig))

	// Register Global Middlewares
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.SanitizeMiddleware())

	// Register Routes
	routes.RegisterAuthRoutes(router)
	routes.RegisterUserRoutes(router)
	routes.RegisterRequestRoutes(router)
	routes.RegisterAdminRoutes(router)

	// Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
