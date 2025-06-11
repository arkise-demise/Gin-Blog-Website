package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var blogpost models.Blog

	// 1. Bind the incoming JSON payload (title, description, image)
	// The 'UserID' field will be 0 initially here, as it's not sent from frontend
	if err := c.ShouldBindJSON(&blogpost); err != nil {
		log.Printf("Error binding payload for CreatePost: %v\n", err.Error()) // Use log.Printf for detailed errors
		c.JSON(400, gin.H{"message": "Invalid payload!"})
		return
	}

	// 2. Get the UserID from the Gin context (set by AuthMiddleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for CreatePost. Is AuthMiddleware properly setting it?")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}

	// Assert type and convert the userID from string to uint
	userIDStr, ok := userIDVal.(string)
	if !ok {
		log.Printf("Error: UserID in context is not a string, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32) // Parse as uint32
	if err != nil {
		log.Printf("Error parsing userID '%s' from context to uint: %v\n", userIDStr, err)
		c.JSON(500, gin.H{"message": "Invalid user ID format."})
		return
	}

	// 3. Assign the extracted UserID to the blogpost
	blogpost.UserID = uint(userID) // Cast to uint

	// 4. (Optional but recommended) Verify if the user exists in the database
	var user models.User
	// Use database.DB directly if it's a globally available connection
	// Or get it from context if you set it there, e.g., db := c.MustGet("db").(*gorm.DB)
	result := database.DB.First(&user, blogpost.UserID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with ID %d not found for post creation.\n", blogpost.UserID)
			c.JSON(400, gin.H{"message": "Associated user not found, cannot create post."})
			return
		}
		log.Printf("Database error checking user existence for ID %d: %v\n", blogpost.UserID, result.Error)
		c.JSON(500, gin.H{"message": "Database error during post creation."})
		return
	}

	// 5. Create the blog post in the database
	if err := database.DB.Create(&blogpost).Error; err != nil {
		log.Printf("Error creating blog post in database: %v\n", err.Error()) // More specific error logging
		// Check for specific database errors if needed, e.g., foreign key violations
		if err.Error() == `pq: insert or update on table "blogs" violates foreign key constraint "fk_blogs_user"` {
			c.JSON(400, gin.H{"message": "Invalid user associated with post (foreign key constraint violated)."})
		} else {
			c.JSON(500, gin.H{"message": "Failed to create post due to database error."})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Congratulations!, your post is done!", "post": blogpost})
}

func GetAllPost(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getblog []models.Blog

	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getblog)
	database.DB.Model(&models.Blog{}).Count(&total)

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(200, gin.H{
		"data": getblog,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"last_page": lastPage,
		},
	})
}

func GetPostById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var blogpost models.Blog
	database.DB.Where("id = ?", id).Preload("User").First(&blogpost)
	c.JSON(200, gin.H{"data": blogpost})
}

func UpdatePostById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	blog := models.Blog{
		Id: uint(id),
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		fmt.Println("Unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid payload!"})
		return
	}
	database.DB.Model(&blog).Updates(blog)
	c.JSON(200, gin.H{"message": "Post updated successfully"})
}

func UniquePost(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := utils.ParseJwt(cookie)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var blogs []models.Blog

	result := database.DB.Where("user_id = ?", id).Preload("User").Find(&blogs)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Could not retrieve blogs"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No blogs found for this user"})
		return
	}

	c.JSON(200, gin.H{"data": blogs})
}

func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID"})
		return
	}

	blog := models.Blog{
		Id: uint(id),
	}

	deleteQuery := database.DB.Delete(&blog)

	if deleteQuery.RowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Oops! Record not found"})
		return
	}

	if deleteQuery.Error != nil {
		c.JSON(500, gin.H{"message": "Error deleting the post"})
		return
	}

	c.JSON(200, gin.H{"message": "Post deleted successfully!"})
}
