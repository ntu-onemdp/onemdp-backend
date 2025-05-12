package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func RetrieveImageHandler(c *gin.Context) {
	utils.Logger.Debug().Interface("header", c.Request.Header).Msg("Received request to retrieve image")

	id := c.Param("id")
	if id == "" {
		utils.Logger.Error().Msg("No ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ID provided"})
		return
	}

	img, err := repositories.Images.Get(id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to retrieve image")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image"})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", img)
}
