package middleware

import (
	"Gin-Blog-Website/models"
	"log"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the authenticated user has the 'admin' role.
// This middleware must be applied *after* AuthMiddleware.
func AdminMiddleware(c *gin.Context) {
	// Retrieve the user object set by AuthMiddleware
	userVal, exists := c.Get("user")
	if !exists {
		log.Println("AdminMiddleware: User object not found in context. AuthMiddleware might not have run or failed.")
		c.AbortWithStatusJSON(500, gin.H{"message": "Server Error: User context missing. Ensure AuthMiddleware runs first."})
		return
	}

	// Type assert the user object
	user, ok := userVal.(models.User)
	if !ok {
		log.Printf("AdminMiddleware: User context is of unexpected type %T.\n", userVal)
		c.AbortWithStatusJSON(500, gin.H{"message": "Server Error: Invalid user context type."})
		return
	}

	// Check if the user's role is 'admin'
	if user.Role != "admin" {
		log.Printf("AdminMiddleware: User %s (ID: %d) attempted to access admin route without 'admin' role (Role: %s).", user.Email, user.Id, user.Role)
		c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden: Admin access required."})
		return
	}

	log.Printf("AdminMiddleware: User %s (ID: %d) successfully granted admin access.", user.Email, user.Id)
	c.Next() // User is an admin, proceed to the handler
}
