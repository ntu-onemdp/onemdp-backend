package files

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func UploadFileHandler(c *gin.Context) {
	utils.Logger.Debug().Msg("Received request to upload file")

	author := services.JwtHandler.GetUidFromJwt(c)

	form, err := c.MultipartForm()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing multipart form")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error parsing multipart form",
			"error":   err,
		})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		utils.Logger.Warn().Msg("No file received.")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No files received",
			"error":   "No files received",
		})
		return
	}

	filegroup := c.PostForm("filegroup")

	for _, file := range files {
		utils.Logger.Debug().Msg("filename: " + file.Filename)

		// Save file metadata to database
		dbFile, err := services.Files.Create(author, file.Filename, &filegroup)
		if err != nil {
			utils.Logger.Error().Err(err).Str("filename", file.Filename).Msg("Failed to save file to database")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error encountered when saving file metadata to database",
				"error":   err,
			})
			return
		}

		// Upload file to GCS
		if err := services.GCSFileServiceInstance.Upload(file, dbFile.GCSFilename); err != nil {
			utils.Logger.Error().Err(err).Msg("Error uploading file to GCS")

			// Revert change in db
			services.Files.Revert(dbFile.FileId)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error encountered when uploading file to GCS",
				"error":   err,
			})
			return
		}

		// Asynchronously upload file to Eduvisor
		go func() {
			if err := services.Eduvisor.Upload(file); err != nil {
				utils.Logger.Warn().Msg("Error uploading file to Eduvisor")
			}
		}()

	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "File uploaded successfully",
	})
}
