package main

import (
	"log"
	"os"

	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models" // <--- IMPORTANT: Import your models package
	"Gin-Blog-Website/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	// --- IMPORTANT: AutoMigrate your new Comment model here ---
	// Ensure User and Blog are also included if they are not migrated elsewhere
	database.DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Comment{})
	// -----------------------------------------------------------

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	routes.Setup(app) // Setup your routes AFTER CORS middleware

	err = app.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
