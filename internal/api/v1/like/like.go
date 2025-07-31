package like

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Handle like content request
func LikeContentHandler(c *gin.Context) {
	// Get user uid from JWT token
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Check if content id is valid and get content ID.
	contentID := services.GetContentID(c)

	// Check if user has already liked content
	if services.Likes.HasLiked(uid, contentID) {
		utils.Logger.Warn().Msgf("User %s has already liked content of id %s", uid, contentID)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has already liked content",
		})
		return
	}

	// Trigger like
	if err := services.Likes.CreateLike(uid, contentID); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error creating like for content id %s and uid %s", contentID, uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error liking content",
			"error":   err.Error(),
		})
		return
	}

	utils.Logger.Trace().Msgf("uid %s has liked content %s successfully", uid, contentID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "content liked",
	})
}

// Handle unlike content request
func UnlikeContentHandler(c *gin.Context) {
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Check content exists and get content ID
	contentID := services.GetContentID(c)

	// Check if user has liked content yet
	if !services.Likes.HasLiked(uid, contentID) {
		utils.Logger.Warn().Msgf("uid %s has not liked content id %s yet", uid, contentID)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has not liked content",
		})
		return
	}

	// Unlike content
	if err := services.Likes.RemoveLike(uid, contentID); err != nil {
		utils.Logger.Err(err).Msg("Error removing like")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error unliking content",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "content unliked successfully",
	})
}
