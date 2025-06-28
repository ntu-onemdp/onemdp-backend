package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func DeletePostsHandler(c *gin.Context) {
	postID := c.Param("post_id")

	// For debugging purposes
	utils.Logger.Info().Str("postID", postID).Msg("Delete post request received")

	// Get author from JWT token
	author := services.JwtHandler.GetUidFromJwt(c)
	utils.Logger.Info().Msg("Delete post request received from " + author)

	// Delete post
	err := services.Posts.DeletePost(postID, author)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting post")

		if (err == utils.ErrUnauthorized{}) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized to delete post. You need to be a staff/admin or the original author to delete the post",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error deleting post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post successfully deleted",
	})
}
