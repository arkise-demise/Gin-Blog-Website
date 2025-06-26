package routes

import (
	"Gin-Blog-Website/controller"
	"Gin-Blog-Website/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// Public Routes - Accessible without authentication
	app.POST("/api/register", controller.RegisterController)
	app.POST("/api/login", controller.LoginController)

	// Public Post & Comment Viewing (ONLY APPROVED CONTENT)
	app.GET("/api/posts", controller.GetAllPost)
	app.GET("/api/posts/:id", controller.GetPostById)
	app.GET("/api/posts/:id/comments", controller.GetCommentsByPostID)

	app.GET("/api/users/:id/profile", controller.GetUserProfile)

	// Authenticated User Routes - Requires AuthMiddleware
	auth := app.Group("/api") // Grouping authenticated routes under /api
	auth.Use(middleware.AuthMiddleware)
	{
		auth.GET("/user", controller.UserGetController)
		auth.POST("/logout", controller.LogoutController)

		// Post-related routes for authenticated users
		auth.POST("/posts", controller.CreatePost)
		auth.GET("/posts/user", controller.GetMyPosts)
		auth.PUT("/posts/:id", controller.UpdatePostById)
		auth.DELETE("/posts/:id", controller.DeletePost)

		// Comment-related route for authenticated users
		auth.POST("/posts/:id/comments", controller.CreateComment)

		auth.GET("/my-profile", controller.GetMyProfile)
		auth.PUT("/my-profile", controller.UpdateMyProfile)

		// File Upload route
		auth.POST("/upload", controller.Upload)
	}

	// Admin Routes - Require both AuthMiddleware AND AdminMiddleware
	admin := app.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	{
		// User Management
		admin.GET("/users", controller.GetAllUsersForAdmin)
		admin.PUT("/users/:id/role", controller.UpdateUserRoleAsAdmin)
		admin.DELETE("/users/:id", controller.DeleteUserAsAdmin)

		// Content Approval - Posts
		admin.GET("/posts/pending", controller.GetPendingPostsForAdmin)
		admin.PUT("/posts/:id/approve", controller.ApprovePostAsAdmin)
		admin.PUT("/posts/:id/reject", controller.RejectPostAsAdmin)

		// Content Approval - Comments
		admin.GET("/comments/pending", controller.GetPendingCommentsForAdmin)
		admin.PUT("/comments/:id/approve", controller.ApproveCommentAsAdmin)
		admin.PUT("/comments/:id/reject", controller.RejectCommentAsAdmin)

		// General Admin Content Moderation (can view/delete any content, regardless of approval)
		admin.GET("/posts", controller.GetAllPostsForAdmin)
		admin.DELETE("/posts/:id", controller.DeletePostAsAdmin)
		admin.GET("/comments", controller.GetAllCommentsForAdmin)
		admin.DELETE("/comments/:id", controller.DeleteCommentAsAdmin)
	}
}
