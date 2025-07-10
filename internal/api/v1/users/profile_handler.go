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
	uid := c.Param("uid")
	if uid == "" {
		utils.Logger.Error().Msg("UID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID is required"})
		return
	}

	utils.Logger.Info().Str("uid", uid).Msg(fmt.Sprintf("Get user profile request received for %s", uid))

	profile, err := services.Users.GetProfile(uid)
	if err != nil {
		utils.Logger.Warn().Msg("profile may be nil. returning 404 here")
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, profile)
}
