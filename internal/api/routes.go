package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/admin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/images"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/posts"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/threads"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/users"
)

/*
################################
||                            ||
||       PUBLIC ROUTES        ||
||                            ||
################################
*/
// Register unprotected login route
func RegisterLoginRoute(router *gin.Engine) {
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		auth.LoginHandler(c)
	})
}

/*
################################
||                            ||
||       STUDENT ROUTES       ||
||                            ||
################################

Routes that are accessible to any authenticated user.
*/
// Image routes
func RegisterImageRoutes(router *gin.RouterGroup) {
	// [AE-88] GET /api/v1/images/:id
	router.GET("/:id", func(c *gin.Context) {
		images.RetrieveImageHandler(c)
	})

	// [AE-87] POST /api/v1/images/upload
	router.POST("/upload", func(c *gin.Context) {
		images.UploadImageHandler(c)
	})
}

// Student routes. Current implementation: jwt verification performed inside handler.
func RegisterStudentUserRoutes(router *gin.RouterGroup) {
	// [AE-6] GET /api/v1/users/:uid
	router.GET("/:uid", func(c *gin.Context) {
		users.GetProfileHandler(c)
	})

}

// Routes starting with /threads
func RegisterThreadRoutes(router *gin.RouterGroup) {
	// [AE-16] POST /api/v1/threads/new
	router.POST("/new", func(c *gin.Context) {
		threads.CreateThreadHandler(c)
	})

	// [AE-25] POST /api/v1/threads/:thread_id/like
	router.POST("/:thread_id/like", func(c *gin.Context) {
		threads.LikeThreadHandler(c)
	})

	// [AE-14] GET /api/v1/threads?size=25&sort=time_created&desc=true&timestamp=0
	router.GET("/", func(c *gin.Context) {
		threads.GetAllThreadsHandler(c)
	})

	// [AE-20] GET /api/v1/threads/:thread_id
	router.GET("/:thread_id", func(c *gin.Context) {
		threads.GetOneThreadHandler(c)
	})

	// [AE-17] DELETE /api/v1/threads/:thread_id
	router.DELETE("/:thread_id", func(c *gin.Context) {
		threads.DeleteThreadHandler(c)
	})

	// [AE-86] DELETE /api/v1/threads/:thread_id/like
	router.DELETE("/:thread_id/like", func(c *gin.Context) {
		threads.UnlikeThreadHandler(c)
	})
}

// Routes starting with /posts
func RegisterPostRoutes(router *gin.RouterGroup) {
	// [AE-21] POST /api/v1/posts/new
	router.POST("/new", func(c *gin.Context) {
		posts.NewPostHandler(c)
	})

	// [AE-89] GET /api/v1/posts/:post_id
	router.GET("/:post_id", func(c *gin.Context) {
		posts.GetPostHandler(c)
	})

	// [AE-26] POST /api/v1/posts/:post_id/like
	router.POST("/:post_id/like", func(c *gin.Context) {
		posts.LikePostHandler(c)
	})

	// [AE-23] POST /api/v1/posts/:post_id/edit
	router.POST("/:post_id/edit", func(c *gin.Context) {
		posts.UpdatePostHandler(c)
	})

	// [AE-24] DELETE /api/v1/posts/:post_id
	router.DELETE("/:post_id", func(c *gin.Context) {
		posts.DeletePostsHandler(c)
	})

	// [AE-85] DELETE /api/v1/posts/:post_id/like
	router.DELETE("/:post_id/like", func(c *gin.Context) {
		posts.UnlikePostHandler(c)
	})
}

/*
################################
||                            ||
||        ADMIN ROUTES        ||
||                            ||
################################
*/
// Register admin routes for user management
func RegisterAdminUserRoutes(router *gin.RouterGroup) {
	// Register create users handler
	router.POST("/users/create", func(c *gin.Context) {
		admin.CreateUsersHandler(c)
	})

	// [AE-9] GET /api/v1/admin/users
	router.GET("/users", func(c *gin.Context) {
		admin.GetAllUsersHandler(c)
	})

	// [AE-8] GET /api/v1/admin/users/:username
	router.GET("/users/:uid", func(c *gin.Context) {
		admin.GetOneUserHandler(c)
	})

	// [AE-12] /api/v1/admin/users/update-role
	router.POST("/users/update-role", func(c *gin.Context) {
		admin.UpdateRoleHandler(c)
	})
}
