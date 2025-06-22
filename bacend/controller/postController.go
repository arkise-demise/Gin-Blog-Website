package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var blogpost models.Blog

	if err := c.ShouldBindJSON(&blogpost); err != nil {
		log.Printf("Error binding payload for CreatePost: %v\n", err.Error())
		c.JSON(400, gin.H{"message": "Invalid payload!"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for CreatePost. AuthMiddleware might not be set up correctly or user not logged in.")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}

	// --- FIX START ---
	// UserID is now stored as uint by AuthMiddleware
	userID, ok := userIDVal.(uint) // Change from string to uint
	if !ok {
		log.Printf("Error: UserID in context is not a uint, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	// No need for strconv.ParseUint here anymore because it's already a uint
	// --- FIX END ---

	blogpost.UserID = userID // Directly assign the uint userID
	// NEW: Set IsApproved to false by default for pending approval
	blogpost.IsApproved = false

	var user models.User
	// Use the uint userID directly for the database query
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

	if err := database.DB.Create(&blogpost).Error; err != nil {
		log.Printf("Error creating blog post in database: %v\n", err.Error())
		if strings.Contains(err.Error(), "foreign key constraint") {
			c.JSON(400, gin.H{"message": "Invalid user associated with post (foreign key constraint violated)."})
		} else {
			c.JSON(500, gin.H{"message": "Failed to create post due to database error."})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Post submitted for approval!", "post": blogpost})
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

	// NEW: Only retrieve approved posts for public view
	query := database.DB.Where("is_approved = ?", true).Preload("User")
	query.Offset(offset).Limit(limit).Find(&getblog)
	query.Model(&models.Blog{}).Count(&total) // Count only approved posts

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}
	var blogpost models.Blog
	// NEW: For a single post, also check if it's approved for public viewing
	result := database.DB.Where("id = ? AND is_approved = ?", id, true).Preload("User").First(&blogpost)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Post not found or not yet approved."})
			return
		}
		log.Printf("Database error fetching post by ID %d: %v\n", id, result.Error)
		c.JSON(500, gin.H{"message": "Database error retrieving post."})
		return
	}
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
		log.Println("Error: UserID not found in context for UpdatePostById. AuthMiddleware might not be set up correctly or user not logged in.")
		c.JSON(500, gin.H{"message": "Authentication context missing."})
		return
	}
	// --- FIX START ---
	currentUserID, ok := userIDVal.(uint) // Change from string to uint
	if !ok {
		log.Printf("Error: UserID in context is not a uint, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	// No need for strconv.ParseUint here anymore
	// --- FIX END ---

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
	if existingPost.UserID != currentUserID { // Directly compare uints
		log.Printf("Unauthorized attempt to update post %d by user %d. Owner is %d.\n", postID, currentUserID, existingPost.UserID)
		c.JSON(403, gin.H{"message": "You are not authorized to update this post."})
		return
	}

	// 5. Bind the incoming JSON payload for updates
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Printf("Error binding update payload for Post ID %d: %v\n", postID, err.Error())
		c.JSON(400, gin.H{"message": "Invalid update payload!"})
		return
	}

	// 6. Update the post in the database
	updateResult := database.DB.Model(&existingPost).Updates(updates)
	if updateResult.Error != nil {
		log.Printf("Error updating post %d in database: %v\n", postID, updateResult.Error)
		c.JSON(500, gin.H{"message": "Failed to update post due to database error."})
		return
	}

	if updateResult.RowsAffected == 0 {
		log.Printf("No changes applied when updating post %d. Possibly same data sent.", postID)
		c.JSON(200, gin.H{"message": "Post updated, but no new changes were applied (possibly same data)."})
		return
	}

	// 7. Respond with success
	c.JSON(200, gin.H{"message": "Post updated successfully!", "post": existingPost})
}

func GetMyPosts(c *gin.Context) { // Renamed from UniquePost for clarity
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for GetMyPosts. AuthMiddleware might not be set up correctly or user not logged in.")
		c.JSON(401, gin.H{"message": "Unauthorized: User ID missing."})
		return
	}

	// --- FIX START ---
	currentUserID, ok := userIDVal.(uint) // Change from string to uint
	if !ok {
		log.Printf("Error: UserID in context is not a uint, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	// --- FIX END ---

	var blogs []models.Blog
	// Use the uint currentUserID directly for the database query
	result := database.DB.Where("user_id = ?", currentUserID).Preload("User").Find(&blogs)
	if result.Error != nil {
		log.Printf("Error retrieving posts for user %d: %v\n", currentUserID, result.Error)
		c.JSON(500, gin.H{"message": "Could not retrieve your posts."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(200, gin.H{"data": []models.Blog{}, "message": "No posts found for this user."})
		return
	}

	c.JSON(200, gin.H{"data": blogs})
}

func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}

	// 1. Get current UserID from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Println("Error: UserID not found in context for DeletePost. AuthMiddleware might not be set up correctly or user not logged in.")
		c.JSON(401, gin.H{"message": "Unauthorized: User ID missing."})
		return
	}

	// --- FIX START ---
	currentUserID, ok := userIDVal.(uint) // Change from string to uint
	if !ok {
		log.Printf("Error: UserID in context is not a uint, got %T\n", userIDVal)
		c.JSON(500, gin.H{"message": "Invalid user ID format in context."})
		return
	}
	// No need for strconv.ParseUint here anymore
	// --- FIX END ---

	// 2. Find the post to be deleted
	var blog models.Blog
	result := database.DB.First(&blog, postID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Post with ID %d not found for deletion.\n", postID)
			c.JSON(404, gin.H{"message": "Post not found."})
			return
		}
		log.Printf("Database error fetching post for deletion (ID %d): %v\n", postID, result.Error)
		c.JSON(500, gin.H{"message": "Database error retrieving post for deletion."})
		return
	}

	// 3. Authorization Check: Ensure the current user owns this post
	if blog.UserID != currentUserID { // Directly compare uints
		log.Printf("Unauthorized attempt to delete post %d by user %d. Owner is %d.\n", postID, currentUserID, blog.UserID)
		c.JSON(403, gin.H{"message": "Forbidden: You are not authorized to delete this post."})
		return
	}

	// 4. Proceed with deletion
	deleteResult := database.DB.Delete(&blog) // Delete the retrieved blog object
	if deleteResult.Error != nil {
		log.Printf("Error deleting post %d from database: %v\n", postID, deleteResult.Error)
		c.JSON(500, gin.H{"message": "Failed to delete the post due to a database error."})
		return
	}

	if deleteResult.RowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Post not found or already deleted."})
		return
	}

	c.JSON(200, gin.H{"message": "Post deleted successfully!"})
}
