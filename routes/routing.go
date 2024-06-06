package routes

import (
	"Gin-Blog-Website/controller"
	"Gin-Blog-Website/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Post("/api/register",controller.RegisterController)
	app.Post("/api/login",controller.LoginController)

	app.Use(middleware.AuthMiddleware)
	app.Post("/api/post",controller.CreatePost)
	app.Get("/api/allpost",controller.GetAllPost)
	app.Get("/api/allpost/:id",controller.GetPostById)
	app.Put("/api/updatepost/:id",controller.UpdatePostById)
	app.Get("/api/uniquepost",controller.UniquePost)
    app.Delete("/api/deletepost/:id", controller.DeletePost)






}