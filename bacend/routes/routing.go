package routes

import (
	"Gin-Blog-Website/controller" // Ensure this imports all your controller functions
	"Gin-Blog-Website/middleware" // Ensure middleware is imported

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

	// NEW: Public User Profile Viewing
	app.GET("/api/users/:id/profile", controller.GetUserProfile) // Get any user's public profile by ID

	// Static file serving for uploads (e.g., images)
	app.Static("/api/uploads", "./uploads")

	// Authenticated User Routes - Requires AuthMiddleware
	auth := app.Group("/api") // Grouping authenticated routes under /api
	auth.Use(middleware.AuthMiddleware)
	{
		auth.GET("/user", controller.UserGetController)
		auth.POST("/logout", controller.LogoutController)

		// Post-related routes for authenticated users
		auth.POST("/posts", controller.CreatePost)        // User creates a post (initially pending approval)
		auth.GET("/posts/user", controller.GetMyPosts)    // User can view all their own posts (including pending)
		auth.PUT("/posts/:id", controller.UpdatePostById) // User can update their own posts
		auth.DELETE("/posts/:id", controller.DeletePost)  // User can delete their own posts

		// Comment-related route for authenticated users
		auth.POST("/posts/:id/comments", controller.CreateComment) // User creates a comment (initially pending approval)

		// NEW: My Profile Management (for the authenticated user themselves)
		auth.GET("/my-profile", controller.GetMyProfile)    // Get the authenticated user's own profile
		auth.PUT("/my-profile", controller.UpdateMyProfile) // Update the authenticated user's own profile

		// File Upload route
		auth.POST("/upload", controller.Upload) // For uploading images (e.g., for posts)
	}

	// Admin Routes - Require both AuthMiddleware AND AdminMiddleware
	admin := app.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware, middleware.AdminMiddleware) // Apply both middlewares in order
	{
		// User Management
		admin.GET("/users", controller.GetAllUsersForAdmin)            // Get all users
		admin.PUT("/users/:id/role", controller.UpdateUserRoleAsAdmin) // Update a user's role (e.g., {"role": "admin"})
		admin.DELETE("/users/:id", controller.DeleteUserAsAdmin)       // Delete a user

		// Content Approval - Posts
		admin.GET("/posts/pending", controller.GetPendingPostsForAdmin) // Get all posts awaiting approval
		admin.PUT("/posts/:id/approve", controller.ApprovePostAsAdmin)  // Approve a specific post
		admin.PUT("/posts/:id/reject", controller.RejectPostAsAdmin)    // Reject (and delete) a specific post

		// Content Approval - Comments
		admin.GET("/comments/pending", controller.GetPendingCommentsForAdmin) // Get all comments awaiting approval
		admin.PUT("/comments/:id/approve", controller.ApproveCommentAsAdmin)  // Approve a specific comment
		admin.PUT("/comments/:id/reject", controller.RejectCommentAsAdmin)    // Reject (and delete) a specific comment

		// General Admin Content Moderation (can view/delete any content, regardless of approval)
		admin.GET("/posts", controller.GetAllPostsForAdmin)            // Get all posts (approved or pending)
		admin.DELETE("/posts/:id", controller.DeletePostAsAdmin)       // Delete any post
		admin.GET("/comments", controller.GetAllCommentsForAdmin)      // Get all comments (approved or pending)
		admin.DELETE("/comments/:id", controller.DeleteCommentAsAdmin) // Delete any comment
	}
}
