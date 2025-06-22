package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"log"
	"strconv"

	// Added for string manipulation if needed for error checks
	// Added for time.Now() if not already there, for timestamps
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Admin User Management ---

// GetAllUsersForAdmin retrieves all users in the system.
// Requires AdminMiddleware.
func GetAllUsersForAdmin(c *gin.Context) {
	var users []models.User
	result := database.DB.Find(&users) // Fetch all users

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Admin: Database error retrieving all users: %v\n", result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve users."})
		return
	}

	// For security, do not return password hashes
	for i := range users {
		users[i].Password = nil // Clear the password hash before sending
	}

	c.JSON(200, gin.H{"data": users})
}

// UpdateUserRoleAsAdmin allows an admin to update another user's role.
// Requires AdminMiddleware.
func UpdateUserRoleAsAdmin(c *gin.Context) {
	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID format."})
		return
	}

	var data struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"message": "Invalid data provided. Role is required."})
		return
	}

	// Validate the role
	if data.Role != "user" && data.Role != "admin" {
		c.JSON(400, gin.H{"message": "Invalid role. Role must be 'user' or 'admin'."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "User not found."})
			return
		}
		log.Printf("Admin: Database error finding user %d for role update: %v\n", targetUserID, err)
		c.JSON(500, gin.H{"message": "Failed to update user role."})
		return
	}

	// Prevent admin from changing their own role (optional, but good practice)
	// You can get current admin's ID from c.Get("userID") if needed.
	// For simplicity, we're not adding that check here, but it's a consideration.

	user.Role = data.Role
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("Admin: Database error updating user %d role to '%s': %v\n", targetUserID, data.Role, err)
		c.JSON(500, gin.H{"message": "Failed to update user role due to database error."})
		return
	}

	user.Password = nil // Clear password before sending response
	c.JSON(200, gin.H{"message": "User role updated successfully!", "user": user})
}

// DeleteUserAsAdmin allows an admin to delete any user by their ID.
// Requires AdminMiddleware.
func DeleteUserAsAdmin(c *gin.Context) {
	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID format."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "User not found."})
			return
		}
		log.Printf("Admin: Database error finding user %d for deletion: %v\n", targetUserID, err)
		c.JSON(500, gin.H{"message": "Failed to delete user."})
		return
	}

	// OPTIONAL: Prevent admin from deleting themselves
	// currentUserID := c.MustGet("userID").(string) // Assuming userID is stored as string
	// if strconv.ParseUint(currentUserID, 10, 32) == targetUserID {
	//     c.JSON(403, gin.H{"message": "Admins cannot delete their own account via this endpoint."})
	//     return
	// }

	// Delete the user
	if err := database.DB.Delete(&user).Error; err != nil {
		log.Printf("Admin: Database error deleting user %d: %v\n", targetUserID, err)
		c.JSON(500, gin.H{"message": "Failed to delete user."})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully!"})
}

// --- Admin Content Approval (Blog Posts) ---

// GetPendingPostsForAdmin retrieves all posts that are not yet approved.
// Requires AdminMiddleware.
func GetPendingPostsForAdmin(c *gin.Context) {
	var posts []models.Blog
	// Fetch posts where IsApproved is false
	result := database.DB.Where("is_approved = ?", false).Order("created_at desc").Preload("User").Find(&posts)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Admin: Database error retrieving pending posts: %v\n", result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve pending posts."})
		return
	}

	c.JSON(200, gin.H{"data": posts})
}

// ApprovePostAsAdmin updates a post's status to approved.
// Requires AdminMiddleware.
func ApprovePostAsAdmin(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}

	var post models.Blog
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Post not found."})
			return
		}
		log.Printf("Admin: Database error finding post %d for approval: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to approve post."})
		return
	}

	post.IsApproved = true // Set to approved
	if err := database.DB.Save(&post).Error; err != nil {
		log.Printf("Admin: Database error approving post %d: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to approve post due to database error."})
		return
	}

	c.JSON(200, gin.H{"message": "Post approved successfully!", "post": post})
}

// RejectPostAsAdmin updates a post's status (e.g., marks it as rejected or deletes it).
// Here, we'll delete it upon rejection, but you could implement a "rejected" status instead.
// Requires AdminMiddleware.
func RejectPostAsAdmin(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}

	var post models.Blog
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Post not found."})
			return
		}
		log.Printf("Admin: Database error finding post %d for rejection: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to reject post."})
		return
	}

	// Option 1: Delete the post upon rejection
	if err := database.DB.Delete(&post).Error; err != nil {
		log.Printf("Admin: Database error deleting post %d upon rejection: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to delete post upon rejection."})
		return
	}
	c.JSON(200, gin.H{"message": "Post rejected and deleted successfully!"})
}

// --- Admin Content Approval (Comments) ---

// GetPendingCommentsForAdmin retrieves all comments that are not yet approved.
// Requires AdminMiddleware.
func GetPendingCommentsForAdmin(c *gin.Context) {
	var comments []models.Comment
	// Fetch comments where IsApproved is false
	result := database.DB.Where("is_approved = ?", false).Order("created_at desc").Preload("User").Preload("Blog").Find(&comments)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Admin: Database error retrieving pending comments: %v\n", result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve pending comments."})
		return
	}

	c.JSON(200, gin.H{"data": comments})
}

// ApproveCommentAsAdmin updates a comment's status to approved.
// Requires AdminMiddleware.
func ApproveCommentAsAdmin(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid comment ID format."})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Comment not found."})
			return
		}
		log.Printf("Admin: Database error finding comment %d for approval: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to approve comment."})
		return
	}

	comment.IsApproved = true // Set to approved
	if err := database.DB.Save(&comment).Error; err != nil {
		log.Printf("Admin: Database error approving comment %d: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to approve comment due to database error."})
		return
	}

	c.JSON(200, gin.H{"message": "Comment approved successfully!", "comment": comment})
}

// RejectCommentAsAdmin updates a comment's status (e.g., marks it as rejected or deletes it).
// Here, we'll delete it upon rejection.
// Requires AdminMiddleware.
func RejectCommentAsAdmin(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid comment ID format."})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Comment not found."})
			return
		}
		log.Printf("Admin: Database error finding comment %d for rejection: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to reject comment."})
		return
	}

	// Delete the comment upon rejection
	if err := database.DB.Delete(&comment).Error; err != nil {
		log.Printf("Admin: Database error deleting comment %d upon rejection: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to delete comment upon rejection."})
		return
	}

	c.JSON(200, gin.H{"message": "Comment rejected and deleted successfully!"})
}

// --- General Admin Content Moderation (Existing functions, kept and enhanced) ---

// GetAllCommentsForAdmin retrieves all comments in the system, regardless of approval status.
// This endpoint requires the AdminMiddleware.
func GetAllCommentsForAdmin(c *gin.Context) {
	var comments []models.Comment
	// Preload User and Blog to get associated data easily for admin review
	result := database.DB.Order("created_at desc").Preload("User").Preload("Blog").Find(&comments)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Admin: Database error retrieving all comments: %v\n", result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve comments."})
		return
	}

	c.JSON(200, gin.H{"data": comments})
}

// DeleteCommentAsAdmin allows an admin to delete any comment by its ID.
// This endpoint requires the AdminMiddleware.
func DeleteCommentAsAdmin(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid comment ID format."})
		return
	}

	var comment models.Comment
	// Find the comment first to ensure it exists
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Comment not found."})
			return
		}
		log.Printf("Admin: Database error finding comment %d for deletion: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to delete comment."})
		return
	}

	// Delete the comment
	if err := database.DB.Delete(&comment).Error; err != nil {
		log.Printf("Admin: Database error deleting comment %d: %v\n", commentID, err)
		c.JSON(500, gin.H{"message": "Failed to delete comment."})
		return
	}

	c.JSON(200, gin.H{"message": "Comment deleted successfully!"})
}

// GetAllPostsForAdmin retrieves all blog posts for admin review, including their authors.
// This endpoint requires the AdminMiddleware.
func GetAllPostsForAdmin(c *gin.Context) {
	var posts []models.Blog
	// Preload the User to show author details for each post
	result := database.DB.Order("created_at desc").Preload("User").Find(&posts)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Admin: Database error retrieving all posts for admin: %v\n", result.Error)
		c.JSON(500, gin.H{"message": "Failed to retrieve posts."})
		return
	}

	c.JSON(200, gin.H{"data": posts})
}

// DeletePostAsAdmin allows an admin to delete any post by its ID.
// This endpoint requires the AdminMiddleware.
func DeletePostAsAdmin(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid post ID format."})
		return
	}

	var post models.Blog
	// Find the post first
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Post not found."})
			return
		}
		log.Printf("Admin: Database error finding post %d for deletion: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to delete post."})
		return
	}

	// Delete the post
	if err := database.DB.Delete(&post).Error; err != nil {
		log.Printf("Admin: Database error deleting post %d: %v\n", postID, err)
		c.JSON(500, gin.H{"message": "Failed to delete post."})
		return
	}

	c.JSON(200, gin.H{"message": "Post deleted successfully!"})
}
