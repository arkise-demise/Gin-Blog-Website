// controllers/user_profile_controller.go
package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/platform/cloudinary"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserProfile - Public view of a user's profile
func GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid user ID format."})
		return
	}

	var user models.User
	// Fetch user, explicitly select fields to be public (exclude password)
	// Preload any related data you want to expose publicly (e.g., their posts)
	result := database.DB.Select("id", "first_name", "last_name", "email", "phone", "role", "bio", "profile_picture_url", "location", "website", "created_at").First(&user, targetUserID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "User not found."})
			return
		}
		log.Printf("Database error fetching user profile for ID %d: %v\n", targetUserID, result.Error)
		c.JSON(500, gin.H{"message": "Database error retrieving user profile."})
		return
	}

	c.JSON(200, user)
}

// GetMyProfile - Authenticated user's own profile
func GetMyProfile(c *gin.Context) {
	// userID is already uint due to AuthMiddleware fix
	userID := c.MustGet("userID").(uint) // Get userID from JWT middleware context

	var user models.User
	// Fetch all fields for the authenticated user's own profile
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "Your profile was not found."})
			return
		}
		log.Printf("Database error fetching my profile for user ID %d: %v\n", userID, result.Error)
		c.JSON(500, gin.H{"message": "Database error retrieving your profile."})
		return
	}

	// Remove password hash before sending to client
	user.Password = nil
	c.JSON(200, user)
}

// UpdateMyProfile - Authenticated user updates their own profile
func UpdateMyProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"message": "User not found."})
			return
		}
		c.JSON(500, gin.H{"message": "Database error retrieving user."})
		return
	}

	// For multipart/form-data, use c.Request.FormFile to get the file,
	// and c.PostForm for other fields.
	// You need to decide if you want to allow JSON updates OR multipart updates.
	// For simplicity, we'll process multipart and extract text fields.
	// If you send JSON, you'd use c.ShouldBindJSON().

	// Handle profile picture upload first
	fileHeader, err := c.FormFile("profile_picture") // "profile_picture" is the name of the input field
	var profilePictureURL string
	if err == nil && fileHeader != nil { // File was uploaded
		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open uploaded file: %v", err)
			c.JSON(500, gin.H{"message": "Failed to open uploaded image."})
			return
		}
		defer file.Close() // Ensure the file is closed

		// Upload to Cloudinary
		url, uploadErr := cloudinary.UploadImage(file)
		if uploadErr != nil {
			log.Printf("Cloudinary upload failed: %v", uploadErr)
			c.JSON(500, gin.H{"message": fmt.Sprintf("Failed to upload profile picture: %v", uploadErr)})
			return
		}
		profilePictureURL = url
		log.Printf("Profile picture uploaded to: %s\n", profilePictureURL)

	} else if err != nil && err.Error() != "http: no such file" {
		// Log other errors besides "no such file" which means no file was sent
		log.Printf("Error getting form file 'profile_picture': %v\n", err)
		c.JSON(400, gin.H{"message": "Error processing profile picture."})
		return
	}

	// Update text fields from form data
	updates := make(map[string]interface{})

	if firstName := c.PostForm("first_name"); firstName != "" {
		updates["FirstName"] = firstName
	}
	if lastName := c.PostForm("last_name"); lastName != "" {
		updates["LastName"] = lastName
	}
	if phone := c.PostForm("phone"); phone != "" {
		updates["Phone"] = phone
	}
	if bio := c.PostForm("bio"); bio != "" {
		updates["Bio"] = bio
	}
	if location := c.PostForm("location"); location != "" {
		updates["Location"] = location
	}
	if website := c.PostForm("website"); website != "" {
		updates["Website"] = website
	}

	// Add the uploaded URL if it exists
	if profilePictureURL != "" {
		updates["ProfilePictureURL"] = profilePictureURL
	} else if clearPic := c.PostForm("clear_profile_picture"); clearPic == "true" {
		// Option to clear the existing profile picture
		updates["ProfilePictureURL"] = ""
	}

	if len(updates) == 0 {
		c.JSON(200, gin.H{"message": "No fields to update or no changes detected."})
		return
	}

	// Update the user in the database
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		log.Printf("Database error updating user profile %d: %v\n", userID, err)
		c.JSON(500, gin.H{"message": "Failed to update profile due to database error."})
		return
	}

	// Fetch the updated user to return the latest state (excluding password)
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error refetching user %d after update: %v\n", userID, err)
		// Still return success, but log the error for refetch
		c.JSON(200, gin.H{"message": "Profile updated successfully, but failed to refetch updated data."})
		return
	}
	user.Password = nil // Ensure password is not sent back
	c.JSON(200, gin.H{"message": "Profile updated successfully!", "user": user})
}
