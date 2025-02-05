package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
)

// Register unprotected login route
func RegisterLoginRoute(router *gin.Engine, handler *auth.LoginHandler) {
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		handler.HandleLogin(c)
	})
}
