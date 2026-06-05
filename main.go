package main

import (
	"log"
	"os"

	"backend/config"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize In-Memory Cache
	config.InitCache()

	// Initialize Database Connection
	// Note: You must have PostgreSQL running and credentials configured via environment variables
	// DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT
	config.InitDB()

	// Initialize Gin Router
	router := gin.Default()

	// Register Routes
	routes.RegisterAuthRoutes(router)

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
