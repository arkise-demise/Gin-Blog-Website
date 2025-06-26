package controller

import (
	"log"
	"net/http"

	
	"Gin-Blog-Website/platform/cloudinary" 

	"github.com/gin-gonic/gin" 
	
)

// Upload handles image uploads to Cloudinary for blog posts.
func Upload(c *gin.Context) {
	// Get the file from the form. We expect a single file input field named "image".
	fileHeader, err := c.FormFile("image") 
	if err != nil {
		log.Printf("Upload Error: Failed to get file from form ('image' field): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve image file. Ensure form field name is 'image'."})
		return
	}

	// Open the uploaded file (multipart.File)
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Failed to open uploaded file for Cloudinary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open image for upload."})
		return
	}
	defer file.Close() // Ensure the file is closed after use

	//  Upload to Cloudinary instead of saving locally ---
	imageURL, uploadErr := cloudinary.UploadImage(file)
	if uploadErr != nil {
		log.Printf("Cloudinary Upload Error: %v\n", uploadErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload image to cloud storage."})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded to cloud successfully!",
		"url":     imageURL, // Return the Cloudinary URL
	})
}

