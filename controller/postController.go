package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func CreatePost(c *fiber.Ctx) error{
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost);err != nil {
		fmt.Println("Unable to parse body")
	}
	if err := database.DB.Create(&blogpost).Error;err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Invalid payload!",
		})
	}
	return c.JSON(fiber.Map{
		"message":"Congratulation!,your post is done!",
	})
}

func GetAllPost(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
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

	return c.JSON(fiber.Map{
		"data": getblog,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": lastPage,
		},
	})
}

func GetPostById(c *fiber.Ctx) error {
	id,_ := strconv.Atoi(c.Params("id"))
	var blogpost models.Blog
	database.DB.Where("id=?",id).Preload("User").First(&blogpost)
	return c.JSON(fiber.Map{
		"data":blogpost,
	})
}


func UpdatePostById(c *fiber.Ctx) error{
	id,_ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id:uint(id),
	}

	if err := c.BodyParser(&blog);err != nil {
		fmt.Println("Unable to parse body")
	}
	database.DB.Model(&blog).Updates(blog)
	return  c.JSON(fiber.Map{
		"message":"post updated successfully",
	})
}

func UniquePost(c *fiber.Ctx) error {
    // Retrieve JWT from cookies
    cookie := c.Cookies("jwt")

    // Parse the JWT to get the user ID
    id, err := utils.ParseJwt(cookie)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unauthorized",
        })
    }

    // Initialize a slice to hold the blog posts
    var blogs []models.Blog

    // Query the database for blogs belonging to the user
    result := database.DB.Where("user_id = ?", id).Preload("User").Find(&blogs)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not retrieve blogs",
        })
    }

    // Check if any blogs were found
    if result.RowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "No blogs found for this user",
        })
    }

    // Return the blog data as JSON
    return c.JSON(fiber.Map{
        "data": blogs,
    })
}


func DeletePost(c *fiber.Ctx) error {
    // Convert the id parameter to an integer
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid post ID",
        })
    }

    // Create a Blog instance with the given ID
    blog := models.Blog{
        Id: uint(id),
    }

    // Delete the blog post from the database
    deleteQuery := database.DB.Delete(&blog)

    // Check if the record was found and deleted
    if deleteQuery.RowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Oops! Record not found",
        })
    }

    // Check for errors during the deletion
    if deleteQuery.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error deleting the post",
        })
    }

    // Respond with a success message
    return c.JSON(fiber.Map{
        "message": "Post deleted successfully!",
    })
}