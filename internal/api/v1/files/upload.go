package files

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func UploadFileHandler(c *gin.Context) {
	utils.Logger.Debug().Interface("header", c.Request.Header).Msg("Received request to upload file")

	author := services.JwtHandler.GetUidFromJwt(c)

	file, err := c.FormFile("file")
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to get file from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received: " + err.Error()})
		return
	}

	filegroup := c.PostForm("filegroup")

	// Save file metadata to database
	if err := services.Files.Create(author, file.Filename, &filegroup); err != nil {
		utils.Logger.Error().Err(err).Str("filename", file.Filename).Msg("Failed to save file to database")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error encountered when saving file metadata to database",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
}
