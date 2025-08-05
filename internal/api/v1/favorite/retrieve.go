package favorite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve list of saved threads
func GetSavedHandler(c *gin.Context) {
	// Get user uid from JWT token
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Get content type
	contentType := c.DefaultQuery("content-type", "threads")

	utils.Logger.Info().Str("uid", uid).Str("content-type", contentType).Msgf("Get saved %s request received.", contentType)

	switch contentType {
	case "threads":
		threads, err := services.Favorites.GetThreads(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error encountered retrieving saved threads",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Saved threads successfully retrieved from database",
			"threads": threads,
		})

	case "articles":
		articles, err := services.Favorites.GetArticles(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error encountered retrieving saved articles",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Saved articles successfully retrieved from database",
			"articles": articles,
		})

	default:
		utils.Logger.Warn().Msgf("Invalid content type: %s", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid content type. Valid types: threads, articles",
		})
	}
}
