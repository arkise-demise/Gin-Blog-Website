package controller

import (
	"Gin-Blog-Website/database"
	"Gin-Blog-Website/models"
	"Gin-Blog-Website/utils"

	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)
	return Re.MatchString(email)
}

func RegisterController(c *gin.Context) {
	var data map[string]interface{}
	var userData models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	// Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		c.JSON(400, gin.H{"message": "Password must be greater than 6 characters!"})
		return
	}

	// Validate email
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.JSON(400, gin.H{"message": "Invalid Email Address!"})
		return
	}

	// Check if email already exists in database
	database.DB.Where("email = ?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		c.JSON(400, gin.H{"message": "Email already exists!"})
		return
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)

	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"user":    user,
		"message": "Account created successfully!",
	})
}

func LoginController(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Unable to parse body")
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		c.JSON(404, gin.H{"message": "Email Address doesn't exist, Please, create an account!"})
		return
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.JSON(400, gin.H{"message": "Incorrect password!"})
		return
	}

	token, err := utils.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.SetCookie("jwt", token, int(24*time.Hour.Seconds()), "/", "", false, true)
	c.JSON(200, gin.H{
		"message": "You have logged in successfully!",
		"user":    user,
	})
}

type Claims struct {
	jwt.StandardClaims
}
