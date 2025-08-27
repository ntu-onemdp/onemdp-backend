package karma

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/karma"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func UpdateKarmaHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Update karma settings request received")

	var req karma.Karma
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Error binding karma request",
		})
		return
	}

	utils.Logger.Debug().Interface("settings", req).Msg("Parsed karma settings from request")

	if err := karma.Service.UpdateSettings(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Error updating karma settings",
		})
		return
	}

	utils.Logger.Info().Msg("Karma settings updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Karma settings updated successfully",
	})
}
