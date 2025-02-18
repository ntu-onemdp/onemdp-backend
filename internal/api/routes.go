package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/admin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
)

// Register unprotected login route
func RegisterLoginRoute(router *gin.Engine, handler *auth.LoginHandler) {
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		handler.HandleLogin(c)
	})
}

// Register create users handler
func RegisterCreateUsersRoute(router *gin.RouterGroup, handler *admin.CreateUserHandler) {
	router.POST("/users/create", func(c *gin.Context) {
		handler.HandleCreateNewUser(c)
	})
}
