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

	// Check if content id is valid
	if !isValidContentID(c) {
		return
	}

	contentID := c.Param("content_id")

	// Check if user has already liked content
	if services.Likes.HasLiked(uid, contentID) {
		utils.Logger.Warn().Msgf("User %s has already liked content of id %s", uid, contentID)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has already liked thread",
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
		"error":   nil,
	})
}

// Handle unlike content request
func UnlikeContentHandler(c *gin.Context) {
	uid := services.JwtHandler.GetUidFromJwt(c)

	contentID := c.Param("content_id")

	// Check content exists
	if !isValidContentID(c) {
		return
	}

	// Check if user has liked content yet
	if !services.Likes.HasLiked(uid, contentID) {
		// TODO
	}
}

// Helper function to check if content exists. If content does not exist, automatically return gin response.
func isValidContentID(c *gin.Context) bool {
	contentID := c.Param("content_id")
	contentType := string(contentID[0]) // first char of content id is the content type

	// Check content exist
	switch contentType {
	case "t":
		if !services.Threads.ThreadExists(contentID) {
			utils.Logger.Warn().Msgf("Thread of thread_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "thread id does not exist",
			})
			return false
		}

	case "p":
		if !services.Posts.PostExists(contentID) {
			utils.Logger.Warn().Msgf("Post of post_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "post id does not exist",
			})
			return false
		}

	case "a":
		if !services.Articles.ArticleExists(contentID) {
			utils.Logger.Warn().Msgf("Article of article_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "article id does not exist",
			})
			return false
		}

	// Invalid content type
	default:
		utils.Logger.Error().Msg("Invalid content type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid content id",
		})
		return false
	}

	return true
}
