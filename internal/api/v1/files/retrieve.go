package files

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func GetFileHandler(c *gin.Context) {
	id := c.Param("file_id")

	utils.Logger.Info().Str("file id", id).Msg("Get file request received.")

	// Get file metadata
	metadata, err := services.Files.GetFilename(id)
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

	// Retrieve file from GCS
	reader, contentType, err := services.GCSFileServiceInstance.Retrieve(metadata.GCSFilename)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("Error retrieving %s from GCS", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error retrieving file from Google Cloud Storage bucket",
			"error":   err,
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+metadata.Filename+"\"")
	c.Status(http.StatusOK)
	c.Header("Content-Type", contentType)

	// Stream the file (helpful for large files)
	io.Copy(c.Writer, reader)
}
