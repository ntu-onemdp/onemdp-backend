package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	utils "github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Verification middleware for non-public routes. Reject if invalid auth token
func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Logger.Trace().Msg("AuthGuard triggered")

		// Retrieve jwt token from auth header
		tokenString := c.Request.Header.Get("Authorization")

		// Remove "Bearer " prefix if included
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claim, err := services.JwtHandler.ParseJwt(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		utils.Logger.Trace().Msg(fmt.Sprintf("claim verified for %s of role %s", claim.Username, claim.Role))
		c.Next()
	}
}

// Verification middleware for admin. Reject if not admin.
func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Logger.Trace().Msg("AdminGuard triggered")

		// Retrieve jwt token from auth header
		tokenString := c.Request.Header.Get("Authorization")

		// Remove "Bearer " prefix if included
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claim, err := services.JwtHandler.ParseJwt(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claim.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. You need to be an admin to access this function."})
			c.Abort()
			return
		}
		utils.Logger.Trace().Msg(fmt.Sprintf("claim verified for %s of role %s", claim.Username, claim.Role))
		c.Next()
	}
}
