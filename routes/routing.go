package routes

import (
	"Gin-Blog-Website/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Post("/api/register",controller.RegisterController)
	app.Post("/api/login",controller.LoginController)

	//app.Use(middleware.AuthMiddleware)
}