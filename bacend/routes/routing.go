package routes

import (
	"Gin-Blog-Website/controller" // Ensure you have the comment controller functions here
	"Gin-Blog-Website/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// Public Routes
	app.POST("/api/register", controller.RegisterController)
	app.POST("/api/login", controller.LoginController)
	app.GET("/api/allpost", controller.GetAllPost)
	app.GET("/api/allpost/:id", controller.GetPostById)
	app.Static("/api/uploads", "./uploads") // Ensure this is always public

	// --- New: Public route to get comments for a post ---
	app.GET("/api/posts/:id/comments", controller.GetCommentsByPostID)

	// Authenticated Routes
	auth := app.Group("/api") // Change this group to /api to simplify paths below
	auth.Use(middleware.AuthMiddleware)
	{
		auth.POST("/post", controller.CreatePost)              // Path is now /api/post
		auth.PUT("/updatepost/:id", controller.UpdatePostById) // Path is now /api/updatepost/:id
		auth.GET("/uniquepost", controller.UniquePost)         // Path is now /api/uniquepost
		auth.DELETE("/deletepost/:id", controller.DeletePost)  // Path is now /api/deletepost/:id
		auth.POST("/upload", controller.Upload)                // Path is now /api/upload
		// auth.Static("uploads", "./uploads") // This should be app.Static in public group.

		auth.POST("/posts/:id/comments", controller.CreateComment) // Path is now /api/posts/:id/comments
	}
}
