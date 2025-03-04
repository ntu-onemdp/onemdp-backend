package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/admin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
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

	router.GET("/", func(ctx *gin.Context) {
		userHandlers.UserProfileHandler.HandleGetUserProfile(ctx)
	})

	router.GET("/password-changed", func(ctx *gin.Context) {
		userHandlers.UserProfileHandler.HandleHasPasswordChanged(ctx)
	})

	router.POST("/change-password", func(c *gin.Context) {
		userHandlers.ChangePasswordHandler.HandleChangeUserPassword(c)
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
