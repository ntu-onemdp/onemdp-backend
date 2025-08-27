package karma

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/karma"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func RetrieveKarmaHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Retrieve karma settings request received")
	settings := karma.Service.GetKarmaSettings()
	if settings == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Karma settings not found",
			"message": "Error retrieving karma settings",
		})
		return
	}

	utils.Logger.Debug().Interface("settings", settings).Msg("Retrieved karma settings")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}
