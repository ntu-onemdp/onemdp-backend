package files

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func GetFileHandler(c *gin.Context) {
	id := c.Param("file_id")

	utils.Logger.Info().Str("file id", id).Msg("Get file request received.")

	// Get GCS filename
	filename, err := services.Files.GetGCSFilename(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": fmt.Sprintf("File %s does not exist", id),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error retrieving Google Cloud Storage filename from postgres",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"filename": filename,
	})
}
