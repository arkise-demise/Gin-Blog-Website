// package controller

// import (
// 	"math/rand"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// var letters = []rune("ashfahashfushfahgagasdf")

// func randLetter(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

// func Upload(c *gin.Context) {
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
// 		return
// 	}
// 	files := form.File["image"]
// 	fileName := ""

// 	for _, file := range files {
// 		fileName = randLetter(5) + "-" + file.Filename
// 		if err := c.SaveUploadedFile(file, "./uploads/"+fileName); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"url": "http://localhost:8080/api/uploads/" + fileName,
// 	})
// }

package controller

import (
	"log"
	"math/rand"
	"net/http" // Make sure this is imported
	"strconv"  // For converting UnixNano to string
	"time"     // For time.Now().UnixNano() and rand.Seed

	"github.com/gin-gonic/gin"
	// Ensure other necessary imports for your post_controller.go are here
	// "Gin-Blog-Website/database"
	// "Gin-Blog-Website/models"
	// "Gin-Blog-Website/utils"
	// "math"
	// "strings"
	// "gorm.io/gorm"
)

// Initialize the random number generator once when the package is loaded.
// This ensures better randomness over multiple calls to randString.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// letters contains a broader set of alphanumeric characters for more robust random string generation.
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// randString generates a random string of length n using the defined letters.
func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Upload handles single file uploads, typically for an "image" associated with a blog post.
func Upload(c *gin.Context) {
	// Get the file from the form. We expect a single file input field named "image".
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Upload Error: Failed to get file from form ('image' field): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve image file. Ensure form field name is 'image'."})
		return
	}

	// Generate a highly unique filename to prevent collisions.
	// This combines a short random string, a nanosecond timestamp, and the original filename.
	uniqueFilename := randString(8) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	filepath := "./uploads/" + uniqueFilename // Path where the file will be saved on the server

	// Save the uploaded file to the specified path.
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		log.Printf("Upload Error: Failed to save uploaded file '%s' to '%s': %v\n", file.Filename, filepath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file on the server."})
		return
	}

	// Construct the URL that the frontend can use to access the uploaded file.
	// This assumes your Gin application serves static files from the `/api/uploads` route.
	fileURL := "http://localhost:8080/api/uploads/" + uniqueFilename

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully!",
		"url":     fileURL, // Return the URL of the saved file
	})
}

// (Keep all your other controller functions like CreatePost, GetAllPost, GetPostById,
// UpdatePostById, UniquePost, DeletePost here, as they were in your previous post_controller.go)
