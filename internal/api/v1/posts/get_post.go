package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve a post by post_id
func GetPostHandler(c *gin.Context) {
	postID := c.Param("post_id")

	// For debugging purposes
	utils.Logger.Info().Str("postID", postID).Msg("Get post request received")

	// Get post
	post, err := services.Posts.GetPost(postID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting post")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error getting post" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"post":    post,
	})
}
