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

	"github.com/gofiber/fiber/v2"
)


func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)
	return Re.MatchString(email)
}

func RegisterController(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data);err != nil {
		fmt.Println("unable to parse body")
	}

	//check if password is less than 6 character

	if len(data["password"].(string))<=6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"password must be greater than 6 character!",
		})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))){
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Invalid Email Address!",
		})
	}

	//check if email already exist in database

	database.DB.Where("email=?",strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id!=0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Email already exist!",
		})
	}

	user := models.User{
		FirstName: data["first_name"].(string),
		LastName: data["last_name"].(string),
		Phone: data["phone"].(string),
		Email: strings.TrimSpace(data["email"].(string)),

	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)

	if err != nil {
		log.Println(err)
	}
	c.Status(200)
	  return c.JSON(fiber.Map{
		"user":user,
		"message":"Account created successfully!",
	  })
}


func LoginController(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data);err!=nil{
		fmt.Println("Unable to parse body")
	}

	var user models.User
	database.DB.Where("email=?",data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message":"Email Address doesn't exist,Please,create an account!",
		})
	}
	if err := user.ComparePassword(data["password"]);err !=nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"incorrect password!",
		})
	}

	token, err := utils.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:"jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}
	c.Cookie((&cookie))
	return c.JSON(fiber.Map{
		"message":"you have logged in successfully!",
		"user":user,
	})
}