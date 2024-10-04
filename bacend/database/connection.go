package database

import (
	"Gin-Blog-Website/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construct DSN for PostgreSQL
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, dbname, port)

	// Connect to the PostgreSQL database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to the PostgreSQL database!")
	} else {
		log.Println("PostgreSQL database connected successfully!")
	}

	// Assign the database connection to the global `DB` variable
	DB = database

	// Automigrate to create tables based on models
	err = database.AutoMigrate(
		&models.User{},
		&models.Blog{},
	)
	if err != nil {
		log.Fatal("Error during AutoMigrate:", err)
	} else {
		log.Println("Database tables migrated successfully!")
	}
}
