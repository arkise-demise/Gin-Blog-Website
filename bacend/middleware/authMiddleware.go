package middleware

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"
	"log"
	"strconv" // <--- NEW: Import strconv for string to uint conversion

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		log.Println("AuthMiddleware: JWT cookie not found or invalid:", err)
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized: Not logged in."})
		return
	}

	// Use utils.ParseJwt directly
	// utils.ParseJwt returns the issuer (which is the user ID as a string) and an error.
	userIDStr, err := utils.ParseJwt(tokenString)
	if err != nil {
		log.Println("AuthMiddleware: Failed to parse token or token is invalid:", err)
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized: Invalid token."})
		return
	}

	// --- NEW: Convert userIDStr to uint ---
	// The base 10 means decimal, 64 means uint64, which is then cast to uint
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Println("AuthMiddleware: Failed to convert userID string to uint:", userIDStr, "Error:", err)
		c.AbortWithStatusJSON(500, gin.H{"message": "Server Error: Invalid User ID format in token."})
		return
	}
	// --- END NEW ---

	var user models.User
	// Fetch the full user object from the database using the converted userID (uint)
	// GORM will now correctly use the uint ID to query the primary key.
	if err := database.DB.Where("id = ?", uint(userID)).First(&user).Error; err != nil { // Ensure it's uint(userID)
		log.Println("AuthMiddleware: User not found from token issuer ID:", userID, "Error:", err)
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized: User not found."})
		return
	}

	// Store both userID (now as uint)
	// AND the full user object (for middlewares like AdminMiddleware)
	c.Set("userID", uint(userID)) // <--- IMPORTANT: Set it as uint here!
	c.Set("user", user)

	// Log the user ID as uint
	log.Printf("AuthMiddleware: User %s (ID: %d, Role: %s) authenticated.", user.Email, user.Id, user.Role)

	c.Next() // Proceed to the next middleware or handler
}
