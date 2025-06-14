package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"
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
	// 1. Get post ID from URL parameter
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		log.Printf("Error parsing post ID from URL: %v\n", err)
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}

	// 2. Get UserID from context (set by AuthMiddleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for UpdatePostById. Is AuthMiddleware properly setting it?")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		log.Printf("Error: UserID in context is not a string, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	currentUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Error parsing currentUserID '%s' from context to uint: %v\n", userIDStr, err)
		c.JSON(500, gin.H{"message": "Invalid user ID format."})
		return
	}

	// 3. Fetch the existing post from the database
	var existingPost models.Blog
	result := database.DB.First(&existingPost, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Post with ID %d not found for update.\n", postID)
			c.JSON(404, gin.H{"message": "Post not found."})
			return
		}
		log.Printf("Database error fetching post for update (ID %d): %v\n", postID, result.Error)
		c.JSON(500, gin.H{"message": "Database error retrieving post."})
		return
	}

	// 4. Authorization Check: Ensure the current user owns the post
	if existingPost.UserID != uint(currentUserID) {
		log.Printf("Unauthorized attempt to update post %d by user %d. Owner is %d.\n", postID, currentUserID, existingPost.UserID)
		c.JSON(403, gin.H{"message": "You are not authorized to update this post."})
		return
	}

	// 5. Bind the incoming JSON payload for updates
	// Use a map to allow partial updates without zeroing out unprovided fields
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Printf("Error binding update payload for Post ID %d: %v\n", postID, err.Error())
		c.JSON(400, gin.H{"message": "Invalid update payload!"})
		return
	}

	// 6. Update the post in the database
	// GORM's Model(&existingPost).Updates(map) will only update provided fields
	updateResult := database.DB.Model(&existingPost).Updates(updates)
	if updateResult.Error != nil {
		log.Printf("Error updating post %d in database: %v\n", postID, updateResult.Error)
		c.JSON(500, gin.H{"message": "Failed to update post due to database error."})
		return
	}

	// If RowsAffected is 0, it means no changes were made (e.g., sent same data or no data)
	if updateResult.RowsAffected == 0 {
		log.Printf("No changes applied when updating post %d. Possibly same data sent.", postID)
		c.JSON(200, gin.H{"message": "Post updated, but no changes were applied (possibly same data)."})
		return
	}

	// 7. Respond with success
	c.JSON(200, gin.H{"message": "Post updated successfully!", "post": existingPost}) // Return updated post if needed
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
		c.JSON(200, gin.H{"data": []models.Blog{}}) 
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
