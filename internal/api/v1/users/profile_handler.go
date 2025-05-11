package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type HasPasswordChangedResponse struct {
	Username         string `json:"username"`
	Password_changed bool   `json:"password_changed"` // False if user has not changed his password yet
}

func HasPasswordChangedHandler(c *gin.Context) {
	// Retrieve username and jwt
	username := c.Param("username")
	tokenString := c.Request.Header.Get("Authorization")

	utils.Logger.Info().Msg(fmt.Sprintf("Password changed query received for %s", username))

	// JWT does not match username
	if !services.JwtHandler.ValidateUsername(username, tokenString) {
		c.JSON(401, nil)
		return
	}

	password_changed, err := services.Users.HasPasswordChanged(username)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(500, nil)
	}

	c.JSON(200, HasPasswordChangedResponse{
		Username:         username,
		Password_changed: password_changed,
	})
}

// Retrieve public user profile information as defined in models.UserProfile
func GetProfileHandler(c *gin.Context) {
	username := c.Param("username")

	utils.Logger.Info().Str("username", username).Msg(fmt.Sprintf("Get user profile request received for %s", username))

	profile, err := services.Users.GetProfile(username)
	if err != nil {
		utils.Logger.Debug().Msg("profile may be nil. returning 404 here")
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(404, nil)
		return
	}

	c.JSON(200, profile)
}
