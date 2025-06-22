package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"

	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time" // Keep time import if used elsewhere, though not directly for User creation in this specific RegisterController

	"github.com/dgrijalva/jwt-go" // Keep if Claims struct is used elsewhere or for clarity
	"github.com/gin-gonic/gin"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)
	return Re.MatchString(email)
}

func RegisterController(c *gin.Context) {
	var data map[string]interface{} // Using interface{} as per your current code
	var userData models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	// Type assertions to ensure data types are correct before access
	password, ok := data["password"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "Password field missing or invalid type."})
		return
	}
	email, ok := data["email"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "Email field missing or invalid type."})
		return
	}
	firstName, ok := data["first_name"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "First name field missing or invalid type."})
		return
	}
	lastName, ok := data["last_name"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "Last name field missing or invalid type."})
		return
	}
	phone, ok := data["phone"].(string) // Phone might be optional, handle accordingly
	if !ok {
		phone = "" // Default empty string if not provided or invalid
	}

	// Check if password is less than 6 characters
	if len(password) <= 6 {
		c.JSON(400, gin.H{"message": "Password must be greater than 6 characters!"})
		return
	}

	// Validate email
	trimmedEmail := strings.TrimSpace(email)
	if !validateEmail(trimmedEmail) {
		c.JSON(400, gin.H{"message": "Invalid Email Address!"})
		return
	}

	// Check if email already exists in database
	database.DB.Where("email = ?", trimmedEmail).First(&userData)
	if userData.Id != 0 { // Check if a user was found
		c.JSON(400, gin.H{"message": "Email already exists!"})
		return
	}

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     trimmedEmail,
		// NEW: Set the default role for new users
		Role: "user",
	}

	user.SetPassword(password)
	err := database.DB.Create(&user).Error // Add .Error to check for database errors

	if err != nil {
		log.Printf("Error creating user in DB: %v\n", err) // Use log.Printf for structured logging
		c.JSON(500, gin.H{"message": "Account creation failed due to server error."})
		return // Return after sending error response
	}

	c.JSON(200, gin.H{
		"user":    user,
		"message": "Account created successfully!",
	})
}

func LoginController(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 { // If user.Id is 0, no user was found
		c.JSON(404, gin.H{"message": "Email Address doesn't exist, Please, create an account!"})
		return
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.JSON(400, gin.H{"message": "Incorrect password!"})
		return
	}

	token, err := utils.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		log.Printf("Error generating JWT: %v\n", err) // Log the error for debugging
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	// Set the cookie
	// domain parameter is empty "" which means it will default to the current domain (localhost in dev)
	// You might want to explicitly set "localhost" for clarity if desired.
	c.SetCookie("jwt", token, int(24*time.Hour.Seconds()), "/", "", false, true)
	c.JSON(200, gin.H{
		"message": "You have logged in successfully!",
		"user":    user, // Returning user data on login might be a security concern depending on fields
	})
}
func UserGetController(c *gin.Context) {
	// The userID and user object are set by AuthMiddleware
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for UserGetController. AuthMiddleware missing or failed.")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}

	// If you store userID as a string in context:
	userIDStr, ok := userIDVal.(string)
	if !ok {
		log.Printf("Error: UserID in context is not a string, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}

	// If you store the full user object in context:
	userVal, exists := c.Get("user")
	if !exists {
		log.Println("Error: User object not found in context for UserGetController.")
		c.JSON(500, gin.H{"message": "User context missing."})
		return
	}
	user, ok := userVal.(models.User)
	if !ok {
		log.Printf("Error: User in context is not of type models.User, got %T\n", userVal)
		// Fallback to fetching from DB if context type assertion fails, using userIDStr
		var fetchedUser models.User
		if err := database.DB.First(&fetchedUser, userIDStr).Error; err != nil {
			log.Printf("Error fetching user from DB with ID %s: %v\n", userIDStr, err)
			c.JSON(500, gin.H{"message": "Failed to retrieve user data."})
			return
		}
		user = fetchedUser // Use the fetched user
	}

	// Important: Do not send password hash to the frontend
	user.Password = nil
	c.JSON(200, gin.H{"user": user})
}

// LogoutController handles user logout by clearing the JWT cookie
func LogoutController(c *gin.Context) {
	// Expire the JWT cookie
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true) // MaxAge -1 immediately expires it

	c.JSON(200, gin.H{"message": "Logout successful!"})
}

type Claims struct { // This struct is not used in the provided functions, but kept if you use it elsewhere
	jwt.StandardClaims
}
