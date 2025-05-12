package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func UploadImageHandler(c *gin.Context) {
	utils.Logger.Debug().Interface("header", c.Request.Header).Msg("Received request to upload image")

	file, err := c.FormFile("file")
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to get file from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received: " + err.Error()})
		return
	}

	// Save the file to the server
	id, err := services.Images.Insert(file)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to save file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}

	utils.Logger.Debug().Str("file_id", id).Msg("File uploaded successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}
