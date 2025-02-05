package auth

import (
	"github.com/gin-gonic/gin"
	utils "github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Verification middleware for admin. Reject if not admin.
func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Logger.Trace().Msg("Guard triggered")

		// Retrieve jwt token from auth header
		tokenString := c.Request.Header.Get("Authorization")

		// Remove "Bearer " prefix if included
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claim, err := utils.ParseJwt(tokenString)

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claim.Role != "admin" {
			c.JSON(401, gin.H{"error": "Unauthorized. You need to be an admin to access this function."})
			c.Abort()
			return
		}

		utils.Logger.Trace().Msg(claim.Username)
		utils.Logger.Trace().Msg(claim.Role)
		c.Next()
	}
}
