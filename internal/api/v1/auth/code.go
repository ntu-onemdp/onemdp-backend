package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/semester"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve enrolment code
func GetCodeHandler(c *gin.Context) {
	code := semester.Service.GetCode()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    code,
	})
}

// Refresh new enrolment code, return the newly generated code to the backend
func RefreshCodeHandler(c *gin.Context) {
	code, err := semester.Service.RefreshCode()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered trying to refresh code")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An error occured trying to refresh enrolment code",
			"error":   err,
		})
		return
	}

	utils.Logger.Info().Msg("Enrolment code successfully refreshed.")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Enrolment code successfully refreshed.",
		"code":    code,
	})
}
