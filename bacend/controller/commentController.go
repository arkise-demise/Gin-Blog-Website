package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"log"
	"strconv"
	"strings"
	"time" // NEW: Import the time package for timestamps

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateComment(c *gin.Context) {
	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid blog post ID format."})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for CreateComment. AuthMiddleware missing or failed.")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		log.Printf("Error: UserID in context is not a string, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Error parsing userID '%s' from context to uint: %v\n", userIDStr, err)
		c.JSON(500, gin.H{"message": "Invalid user ID format."})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding comment payload: %v\n", err.Error())
		c.JSON(400, gin.H{"message": "Comment content is required."})
		return
	}

	comment := models.Comment{
		Content:   input.Content,
		UserID:    uint(userID),
		BlogID:    uint(blogID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		// NEW: Set IsApproved to false by default for pending approval
		IsApproved: false,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		log.Printf("Error creating comment in database: %v\n", err.Error())
		if strings.Contains(err.Error(), "foreign key constraint") {
			c.JSON(400, gin.H{"message": "Invalid user or blog post ID (foreign key error)."})
			return
		}
		c.JSON(500, gin.H{"message": "Failed to create comment due to database error."})
		return
	}

	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(201, gin.H{"message": "Comment submitted for approval!", "comment": comment}) // Changed message
}

func GetCommentsByPostID(c *gin.Context) {
	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid blog post ID format."})
		return
	}

	var comments []models.Comment
	// NEW: Only retrieve approved comments for public view
	result := database.DB.Where("blog_id = ? AND is_approved = ?", blogID, true).Order("created_at asc").Preload("User").Find(&comments)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Database error retrieving comments for blog %d: %v\n", blogID, result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve comments."})
		return
	}
	c.JSON(200, gin.H{"data": comments})
}
