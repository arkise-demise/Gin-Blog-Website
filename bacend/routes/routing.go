// route.go
package routes

import (
	"Gin-Blog-Website/controller"
	"Gin-Blog-Website/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// Public Routes
	app.POST("/api/register", controller.RegisterController)
	app.POST("/api/login", controller.LoginController)
	app.GET("/api/allpost", controller.GetAllPost)
	app.GET("/api/allpost/:id", controller.GetPostById)
	app.Static("api/uploads", "./uploads") // <-- This was also inside the auth group before, but it should be public.

	// Authenticated Routes
	auth := app.Group("/") // Group starts at the root path '/'
	auth.Use(middleware.AuthMiddleware)
	{
		auth.POST("/api/post", controller.CreatePost)
		auth.PUT("/api/updatepost/:id", controller.UpdatePostById)
		auth.GET("/api/uniquepost", controller.UniquePost)
		auth.DELETE("/api/deletepost/:id", controller.DeletePost)
		auth.POST("api/upload", controller.Upload) // <-- Here's the route!
		// auth.Static("api/uploads","./uploads") // If this is inside the auth group, it's problematic for serving uploaded images publicly.
	}
}
