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

	utils.Logger.Info().Str("username", username).Msg(fmt.Sprintf("Get user profile request received for %s", username))

	profile, err := services.Users.GetProfile(username)
	if err != nil {
		utils.Logger.Debug().Msg("profile may be nil. returning 404 here")
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, profile)
}
