// controllers/user_profile_controller.go (create this file)
package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"net/http"

	// Adjust YOUR_MODULE_PATH
	"github.com/gin-gonic/gin"
)

// DTO for updating user profile
type UpdateUserProfileInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Bio               string `json:"bio"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Location          string `json:"location"`
	Website           string `json:"website"`
}

// GetUserProfile retrieves a user's public profile by ID
func GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.Select("id", "first_name", "last_name", "bio", "profile_picture_url", "location", "website").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetMyProfile retrieves the authenticated user's own profile
func GetMyProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Get userID from JWT middleware context
	var user models.User
	// Select all fields for the authenticated user, or specific ones you want to expose
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"}) // Should not happen if auth works
		return
	}
	// Do not return password hash
	user.Password = []byte{}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateMyProfile allows the authenticated user to update their own profile
func UpdateMyProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Get userID from JWT middleware context

	var input UpdateUserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Update fields. Use GORM's Updates method.
	// You might want to add validation here before updating
	database.DB.Model(&user).Updates(models.User{
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Bio:               input.Bio,
		ProfilePictureURL: input.ProfilePictureURL,
		Location:          input.Location,
		Website:           input.Website,
	})

	// Do not return password hash
	user.Password = []byte{}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully!", "data": user})
}
