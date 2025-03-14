package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/admin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/posts"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/threads"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/users"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

/*
################################
||                            ||
||       PUBLIC ROUTES        ||
||                            ||
################################
*/
// Register unprotected login route
func RegisterLoginRoute(router *gin.Engine, handler *auth.LoginHandler) {
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		handler.HandleLogin(c)
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
// Student routes. Current implementation: jwt verification performed inside handler.
func RegisterStudentUserRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	userHandlers := users.InitUserHandlers(db)

	router.GET("/", func(c *gin.Context) {
		userHandlers.UserProfileHandler.HandleGetUserProfile(c)
	})

	// [AE-7] GET /api/v1/users/:username/password-changed
	router.GET("/password-changed", func(c *gin.Context) {
		userHandlers.UserProfileHandler.HandleHasPasswordChanged(c)
	})

}

// Register change password routes
func RegisterChangePasswordRoute(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repositories.AuthRepository{Db: db}
	usersRepo := repositories.UsersRepository{Db: db}
	authService := services.AuthService{AuthRepo: &authRepo, UsersRepo: &usersRepo}
	changePasswordHandler := auth.ChangePasswordHandler{AuthService: &authService}

	// [AE-4] POST /api/v1/auth/:username/change-password
	router.POST("/api/v1/auth/:username/change-password", func(c *gin.Context) {
		changePasswordHandler.HandleChangeUserPassword(c)
	})
}

// Routes starting with /threads
func RegisterThreadRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	threadHandlers := threads.InitThreadHandlers(db)

	// [AE-16] POST /api/v1/threads/new
	router.POST("/new", func(c *gin.Context) {
		threadHandlers.NewThreadHandler.HandleNewThread(c)
	})

	// [AE-25] POST /api/v1/threads/:thread_id/like
	router.POST("/:thread_id/like", func(c *gin.Context) {
		threadHandlers.LikeThreadHandlers.HandleLikeThread(c)
	})

	// [AE-20] GET /api/v1/threads/:thread_id
	router.GET("/:thread_id", func(c *gin.Context) {
		threadHandlers.GetThreadHandler.HandleGetThread(c)
	})

	// [AE-17] DELETE /api/v1/threads/:thread_id
	router.DELETE("/:thread_id", func(c *gin.Context) {
		threadHandlers.DeleteThreadHandler.HandleDeleteThread(c)
	})

	// [AE-86] DELETE /api/v1/threads/:thread_id/like
	// router.DELETE("/:thread_id/like", func(c *gin.Context) {
	// 	threadHandlers.LikeThreadHandlers.HandleUnlikeThread(c)
	// })
}

// Routes starting with /posts
func RegisterPostRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	postHandlers := posts.InitPostHandlers(db)

	// [AE-21] POST /api/v1/posts/new
	router.POST("/new", func(c *gin.Context) {
		postHandlers.NewPostHandler.HandleNewPost(c)
	})

	// [AE-26] POST /api/v1/posts/:post_id/like
	router.POST("/:post_id/like", func(c *gin.Context) {
		postHandlers.LikePostHandlers.HandleLikePost(c)
	})

	// [AE-23] POST /api/v1/posts/:post_id/edit
	router.POST("/:post_id/edit", func(c *gin.Context) {
		postHandlers.UpdatePostHandler.HandleUpdatePost(c)
	})

	// [AE-24] DELETE /api/v1/posts/:post_id
	router.DELETE("/:post_id", func(c *gin.Context) {
		postHandlers.DeletePostHandler.HandleDeletePost(c)
	})

	// [AE-85] DELETE /api/v1/posts/:post_id/like
	router.DELETE("/:post_id/like", func(c *gin.Context) {
		postHandlers.LikePostHandlers.HandleUnlikePost(c)
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
func RegisterAdminUserRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	userHandlers := admin.InitUserHandlers(db)

	// Register create users handler
	router.POST("/users/create", func(c *gin.Context) {
		userHandlers.CreateUserHandler.HandleCreateNewUser(c)
	})

	// [AE-9] GET /api/v1/admin/users
	router.GET("/users", func(c *gin.Context) {
		userHandlers.GetUsersHandler.HandleGetUsers(c)
	})

	// [AE-8] GET /api/v1/admin/users/:username
	router.GET("/users/:username", func(c *gin.Context) {
		userHandlers.GetUsersHandler.HandleGetUser(c)
	})

	// [AE-12] /api/v1/admin/users/update-role
	router.POST("/users/update-role", func(c *gin.Context) {
		userHandlers.UpdateUsersRoleHandler.HandleUpdateUsersRole(c)
	})
}
