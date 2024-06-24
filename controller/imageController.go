package controller

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

var letters = []rune("ashfahashfushfahgagasdf")

func randLetter(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}
	files := form.File["image"]
	fileName := ""

	for _, file := range files {
		fileName = randLetter(5) + "-" + file.Filename
		if err := c.SaveUploadedFile(file, "./uploads/"+fileName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"url": "http://localhost:8080/api/upload/" + fileName,
	})
}
