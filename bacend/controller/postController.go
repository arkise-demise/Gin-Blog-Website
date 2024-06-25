package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var blogpost models.Blog
	if err := c.ShouldBindJSON(&blogpost); err != nil {
		fmt.Println("Unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid payload!"})
		return
	}
	if err := database.DB.Create(&blogpost).Error; err != nil {
		c.JSON(400, gin.H{"message": "Invalid payload!"})
		return
	}
	c.JSON(200, gin.H{"message": "Congratulations!, your post is done!"})
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
