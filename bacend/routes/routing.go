package routes

import (
	"Gin-Blog-Website/controller"
	"Gin-Blog-Website/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	app.POST("/api/register", controller.RegisterController)
	app.POST("/api/login", controller.LoginController)

	auth := app.Group("/")
	auth.Use(middleware.AuthMiddleware)
	{
		auth.POST("/api/post", controller.CreatePost)
		auth.GET("/api/allpost", controller.GetAllPost)
		auth.GET("/api/allpost/:id", controller.GetPostById)
		auth.PUT("/api/updatepost/:id", controller.UpdatePostById)
		auth.GET("/api/uniquepost", controller.UniquePost)
		auth.DELETE("/api/deletepost/:id", controller.DeletePost)
		auth.POST("api/upload",controller.Upload)
		auth.Static("api/uploads","./uploads")
	}
}
