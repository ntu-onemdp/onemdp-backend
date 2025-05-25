package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Handle like post request
func LikePostHandler(c *gin.Context) {
	// Get username from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := services.JwtHandler.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Error parsing JWT token",
			"error":   err.Error(),
		})
		return
	}

	username := claim.Username
	utils.Logger.Info().Msg("Like post request received from " + username)

	// Get post id from URL
	postID := c.Param("post_id")
	utils.Logger.Trace().Str("post_id", postID).Msg("")

	// Check if post exists
	if !services.Posts.PostExists(postID) {
		utils.Logger.Error().Msg("Post does not exist")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Post does not exist",
			"error":   nil,
		})
		return
	}

	// Check if user has liked the post
	hasLiked := services.Likes.HasLiked(username, postID)

	// User has already liked the post, do nothing
	if hasLiked {
		utils.Logger.Trace().Msg("User has already liked post")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has already liked post",
			"error":   nil,
		})
		return
	}

	err = services.Likes.CreateLike(username, postID)
	if err != nil {
		utils.Logger.Err(err).Msg("Error creating like")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error liking post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post liked successfully",
		"error":   nil,
	})
}

// Handle unlike post request
func UnlikePostHandler(c *gin.Context) {
	// Get username from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := services.JwtHandler.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Error parsing JWT token",
			"error":   err.Error(),
		})
		return
	}

	username := claim.Username
	utils.Logger.Info().Msg("Unlike post request received from " + username)

	// Get post id from URL
	postID := c.Param("post_id")
	utils.Logger.Trace().Str("post_id", postID).Msg("")

	// If user has not liked the post, do nothing
	hasLiked := services.Likes.HasLiked(username, postID)
	if !hasLiked {
		utils.Logger.Trace().Msg("User has not liked post")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has not liked post",
			"error":   nil,
		})
		return
	}

	err = services.Likes.RemoveLike(username, postID)
	if err != nil {
		utils.Logger.Err(err).Msg("Error removing like")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error unliking post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post unliked successfully",
		"error":   nil,
	})
}
