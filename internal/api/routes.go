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

	router.POST("/api/v1/auth/:username/change-password", func(c *gin.Context) {
		changePasswordHandler.HandleChangeUserPassword(c)
	})
}

// Routes starting with /threads
func RegisterThreadRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	threadHandlers := threads.InitThreadHandlers(db)

	router.POST("/new", func(c *gin.Context) {
		threadHandlers.NewThreadHandler.HandleNewThread(c)
	})

	// [AE-20] GET /api/v1/threads/:thread_id
	router.GET("/:thread_id", func(c *gin.Context) {
		threadHandlers.GetThreadHandler.HandleGetThread(c)
	})

	router.DELETE("/:thread_id", func(c *gin.Context) {
		threadHandlers.DeleteThreadHandler.HandleDeleteThread(c)
	})
}

// Routes starting with /posts
func RegisterPostRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	postHandlers := posts.InitPostHandlers(db)

	router.POST("/new", func(c *gin.Context) {
		postHandlers.NewPostHandler.HandleNewPost(c)
	})

	router.DELETE("/:postId", func(c *gin.Context) {
		postHandlers.DeletePostHandler.HandleDeletePost(c)
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

	// Register get users handler
	router.GET("/users", func(c *gin.Context) {
		userHandlers.GetUsersHandler.HandleGetUsers(c)
	})

	// Register get individual user handler
	router.GET("/users/:username", func(c *gin.Context) {
		userHandlers.GetUsersHandler.HandleGetUser(c)
	})

	// Register update users role handler
	router.POST("/users/update-role", func(c *gin.Context) {
		userHandlers.UpdateUsersRoleHandler.HandleUpdateUsersRole(c)
	})
}
