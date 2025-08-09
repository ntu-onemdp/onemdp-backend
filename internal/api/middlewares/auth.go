package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	utils "github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

var API_KEY *string

func init() {
	key, found := os.LookupEnv("API_KEY")
	if !found {
		utils.Logger.Warn().Msg("API KEY not set in env variable.")
	}

	utils.Logger.Info().Msg("API key loaded. Use the API key (found in .env) to access the backend without JWT tokens.")
	API_KEY = &key
}

// Verification middleware for non-public routes. Reject if invalid auth token.
// User role is the mininmum level of authorization.
// If the user role is student, it does not perform any database query.
func AuthGuard(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Logger.Trace().Msg("AuthGuard triggered")

		// Retrieve api key from header
		apiKey := c.Request.Header.Get("x-api-key")
		if apiKey != "" {
			// Match API key
			if apiKey == *API_KEY {
				utils.Logger.Info().Msg("Access via API key granted")
				c.Next()
			} else {
				// If this code is reached, there is an unauthorized attempt to access the backend as user requests use JWT and not API key
				utils.Logger.Warn().Msg("Incorrect API key.")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid key"})
				c.Abort()
				return
			}
		}

		// Retrieve jwt token from auth header
		tokenString := c.Request.Header.Get("Authorization")

		// Remove "Bearer " prefix if included
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claim, err := services.JwtHandler.ParseJwt(tokenString)

		if err != nil {
			utils.Logger.Warn().Msg("Invalid token, rejecting claim.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		utils.Logger.Trace().Msgf("Claim verified for %s", claim.Uid)

		// Pass request if min role is student
		if role <= models.Student {
			c.Next()
			return
		}

		utils.Logger.Trace().Str("min role", role.String()).Msg("AdminGuard triggered")

		userRole, err := services.Users.GetRole(claim.Uid)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error fetching user role")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		utils.Logger.Trace().Str("user role", userRole.String()).Msg("User's role fetched from database")

		// Insufficient permission
		if userRole < role {
			utils.Logger.Warn().Str("uid", claim.Uid).Str("user role", userRole.String()).Str("min role", role.String()).Msg("Request rejected, user does not have sufficient permissions")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permissions to access this resource",
			})
			c.Abort()
			return
		}

		utils.Logger.Trace().Msgf("AdminGuard approved for user %s", claim.Uid)
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

		user, err := services.Users.GetProfile(claim.Uid)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error retrieving user profile")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. You need to be an admin to access this function."})
			c.Abort()
			return
		}

		utils.Logger.Trace().Msg(fmt.Sprintf("claim verified for %s of role %s", claim.Uid, user.Role))
		c.Next()
	}
}
