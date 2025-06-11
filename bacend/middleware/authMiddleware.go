package middleware

import (
	"Gin-Blog-Website/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized! (No JWT cookie)"})
		c.Abort()
		return
	}

	// ParseJwt now returns the issuer (user ID) as a string
	userIDStr, err := utils.ParseJwt(cookie)
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized! (Invalid JWT)"})
		c.Abort()
		return
	}

	// Set the user ID in the context for subsequent handlers
	c.Set("userID", userIDStr) // Store the user ID as a string in the context
	c.Next()
}
