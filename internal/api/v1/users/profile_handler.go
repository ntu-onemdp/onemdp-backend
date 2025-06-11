package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve public user profile information as defined in models.UserProfile
func GetProfileHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.Logger.Error().Msg("Username is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	utils.Logger.Info().Str("username", username).Msg(fmt.Sprintf("Get user profile request received for %s", username))
	email := username + "@e.ntu.edu.sg"

	profile, err := services.Users.GetProfile(email)
	if err != nil {
		utils.Logger.Debug().Msg("profile may be nil. returning 404 here")
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, profile)
}
