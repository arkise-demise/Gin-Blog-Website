// main.go
package main

import (
	"log"
	"os"

	"Gin-Blog-Website/database"
	"Gin-Blog-Website/routes"

	"github.com/gin-contrib/cors" // Import the cors middleware
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	app := gin.Default()

	// --- Add CORS Middleware here ---
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                             // Allow your Next.js frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},           // Allow all methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Important for headers sent by frontend
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,  // Allow cookies (for JWT)
		MaxAge:           86400, // Preflight cache for 24 hours
	}))
	// --- End CORS Middleware ---

	routes.Setup(app) // Setup your routes AFTER CORS middleware

	err = app.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
