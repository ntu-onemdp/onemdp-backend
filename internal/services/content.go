package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve content id from request. Returns response if not found or invalid.
func GetContentID(c *gin.Context) string {
	contentID := c.Param("content_id")
	contentType := string(contentID[0]) // first char of content id is the content type

	// Check content exist
	switch contentType {
	case "t":
		if !Threads.ThreadExists(contentID) {
			utils.Logger.Warn().Msgf("Thread of thread_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "thread id does not exist",
			})
			return ""
		}

	case "p":
		if !Posts.PostExists(contentID) {
			utils.Logger.Warn().Msgf("Post of post_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "post id does not exist",
			})
			return ""
		}

	case "a":
		if !Articles.ArticleExists(contentID) {
			utils.Logger.Warn().Msgf("Article of article_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "article id does not exist",
			})
			return ""
		}

	case "c":
		if !Comments.CommentExists(contentID) {
			utils.Logger.Warn().Msgf("Comment of comment_id %s does not exist", contentID)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "comment id does not exist",
			})
			return ""
		}

	// Invalid content type
	default:
		utils.Logger.Error().Msg("Invalid content type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid content id",
		})
		return ""
	}

	return contentID
}
