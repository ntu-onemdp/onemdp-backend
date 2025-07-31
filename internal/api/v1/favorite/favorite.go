package favorite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Handle favorite content request
func FavoriteContentHandler(c *gin.Context) {
	// Get user uid from JWT token
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Check if content id is valid and get content ID.
	contentID := services.GetContentID(c)

	// Check if user has already favorited content
	if services.Favorites.Exists(uid, contentID) {
		utils.Logger.Warn().Msgf("User %s has already favorited content of id %s", uid, contentID)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has already favorited content",
		})
		return
	}

	// Add record to database
	if err := services.Favorites.CreateFavorite(uid, contentID); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error favoriting content id %s and uid %s", contentID, uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error adding to favorites",
			"error":   err.Error(),
		})
		return
	}

	utils.Logger.Debug().Msgf("uid %s has favorited content %s successfully", uid, contentID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "content added to favorites",
	})
}

func RemoveFavoriteHandler(c *gin.Context) {
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Check content exists and get content ID
	contentID := services.GetContentID(c)

	// Check if user has favorited content yet
	if !services.Favorites.Exists(uid, contentID) {
		utils.Logger.Warn().Msgf("uid %s has not favorited content id %s yet", uid, contentID)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has not favorited content",
		})
		return
	}

	// Remove content from favorites
	if err := services.Favorites.RemoveFavorite(uid, contentID); err != nil {
		utils.Logger.Err(err).Msg("Error removing content from favorites")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error removing content from favorites",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "content removed from favorites successfully",
	})
}
