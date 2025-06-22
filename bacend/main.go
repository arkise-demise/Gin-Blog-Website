package main

import (
	"log"
	"os"

	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models" // IMPORTANT: Import your models package
	"Gin-Blog-Website/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Please create one.")
	}

	// Connect to the database
	database.Connect()

	// AutoMigrate all your models to ensure database tables are up-to-date
	// This is crucial for adding the new 'is_approved' columns to 'blogs' and 'comments' tables.
	database.DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Comment{})
	log.Println("Database migrations completed.")

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set. Please set it in .env or system.")
	}

	// Initialize Gin default router
	app := gin.Default()

	// Configure CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // Cache preflight requests for 24 hours
	}))

	// Setup all API routes
	routes.Setup(app)

	// Run the Gin server
	log.Printf("Server starting on :%s\n", port)
	err = app.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
