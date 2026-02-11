package main

import (
	"log"
	"os"

	auth_services "lem-be/auth/services"
	"lem-be/database"
	"lem-be/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize MongoDB connection
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Bootstrap superuser
	log.Println("Bootstrapping superuser...")	
	if err := auth_services.InitSuperuser(database.GetDB()); err != nil {
		log.Printf("Error bootstrapping superuser: %v", err)
	}
	log.Println("Superuser bootstrapped successfully")

	// Initialize Gin router
	r := router.Setup()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
