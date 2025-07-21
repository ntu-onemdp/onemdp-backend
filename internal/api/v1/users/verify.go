package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func VerifyAdminHanlder(c *gin.Context) {
	// Get uid from jwt
	uid := services.JwtHandler.GetUidFromJwt(c)

	hasAdminPermission, err := services.Users.HasAdminPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("Error getting permission for user %s", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	utils.Logger.Info().Bool("has admin permission", hasAdminPermission).Msgf("Request received from user %s to verify admin status.")
	c.JSON(http.StatusOK, gin.H{
		"hasAdminPermission": hasAdminPermission,
	})
}
