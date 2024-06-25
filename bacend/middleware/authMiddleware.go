package middleware

import (
	"Gin-Blog-Website/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		c.Abort()
		return
	}

	if _, err := utils.ParseJwt(cookie); err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized!"})
		c.Abort()
		return
	}

	c.Next()
}
