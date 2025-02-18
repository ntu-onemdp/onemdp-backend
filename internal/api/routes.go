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
// Student routes that involve only user information
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
