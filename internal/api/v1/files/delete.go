package files

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func DeleteFileHandler(c *gin.Context) {
	// Retrieve params
	uid := services.JwtHandler.GetUidFromJwt(c)
	fileID := c.Param("file_id")

	utils.Logger.Info().Str("uid", uid).Str("file ID", fileID).Msgf("Delete file request received from %s for file %s", uid, fileID)

	// Mark file as deleted
	if err := services.Files.Remove(fileID, uid); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error deleting file in database",
			"error":   err,
		})
		return
	}

	// Note: we do not delete the file in GCS for now.

	utils.Logger.Info().Str("uid", uid).Str("file ID", fileID).Msgf("Successfully deleted file %s", fileID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("File %s successfully deleted.", fileID),
	})
}
