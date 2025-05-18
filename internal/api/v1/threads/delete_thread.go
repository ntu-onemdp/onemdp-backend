package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func DeleteThreadHandler(c *gin.Context) {
	threadId := c.Param("thread_id")

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := services.JwtHandler.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	utils.Logger.Info().Msg("Delete thread request received from " + claim.Username)

	err = services.Threads.DeleteThread(threadId, claim)
	if err == utils.NewErrUnauthorized() {
		utils.Logger.Error().Err(err).Msg("User is student and not author. Unauthorized to delete thread")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized to delete thread. You need to be a staff/admin or the original author to delete the thread",
		})
		return
	} else if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Error deleting thread: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Thread deleted successfully",
	})
}
