package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve public user profile information as defined in models.UserProfile
func GetProfileHandler(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		utils.Logger.Warn().Msg("UID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID is required"})
		return
	}

	utils.Logger.Info().Str("uid", uid).Msgf("Get user profile request received for uid %s", uid)

	profile, err := services.Users.GetProfile(uid)
	if err != nil {
		utils.Logger.Warn().Err(err).Msg("profile may be nil. returning 404 here")
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Retrieve public profile photo
func GetProfilePhotoHandler(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		utils.Logger.Warn().Msg("UID is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UID field is empty",
		})
		return
	}

	utils.Logger.Info().Msgf("Get user profile photo request received for uid %s", uid)

	photo, err := services.Users.GetProfilePhoto(uid)
	if err != nil {
		utils.Logger.Warn().Err(err).Msgf("User %s not found.", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user profile not found",
		})
		return
	}

	// Profile photo not found
	if photo == nil {
		utils.Logger.Warn().Msgf("User %s has no profile photo set", uid)
		c.JSON(http.StatusNoContent, gin.H{
			"message": "image not found",
		})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", photo)
}
