package api

import (
	"github.com/gin-gonic/gin"
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
*/
// Student routes. Current implementation: jwt verification performed inside handler.
func RegisterStudentUserRoutes(router *gin.RouterGroup, handler *users.ProfileHandler) {
	router.GET("/", func(ctx *gin.Context) {
		handler.HandleGetUserProfile(ctx)
	})

	router.GET("/password-changed", func(ctx *gin.Context) {
		handler.HandleHasPasswordChanged(ctx)
	})
}

/*
################################
||                            ||
||        ADMIN ROUTES        ||
||                            ||
################################
*/
// Register create users handler
func RegisterCreateUsersRoute(router *gin.RouterGroup, handler *admin.CreateUserHandler) {
	router.POST("/users/create", func(c *gin.Context) {
		handler.HandleCreateNewUser(c)
	})
}

// Register get users handler
func RegisterGetUsersRoutes(router *gin.RouterGroup, handler *admin.GetUsersHandler) {
	router.GET("/users", func(c *gin.Context) {
		handler.HandleGetUsers(c)
	})

	router.GET("/users/:username", func(c *gin.Context) {
		handler.HandleGetUser(c)
	})
}

// Register update users role handler
func RegisterUpdateUserRoleRoute(router *gin.RouterGroup, handler *admin.UpdateUsersRoleHandler) {
	router.POST("/users/update-role", func(c *gin.Context) {
		handler.HandleUpdateUsersRole(c)
	})
}
