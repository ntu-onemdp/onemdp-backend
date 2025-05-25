package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Handle like thread request
func LikeThreadHandler(c *gin.Context) {
	// Get username form JWT token
	username := services.JwtHandler.GetUsernameFromJwt(c)
	if username == "" {
		return
	}

	// Check thread exists
	threadID := c.Param("thread_id")
	utils.Logger.Trace().Str("thread_id", threadID).Msg("")
	if !services.Threads.ThreadExists(threadID) {
		utils.Logger.Error().Msg("Thread does not exist")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Thread does not exist",
			"error":   nil,
		})
	}

	// Check if user has liked the thread
	hasLiked := services.Likes.HasLiked(username, threadID)
	if hasLiked {
		utils.Logger.Trace().Msg("User has already liked thread")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has already liked thread",
			"error":   nil,
		})
		return
	}

	err := services.Likes.CreateLike(username, threadID)
	if err != nil {
		utils.Logger.Err(err).Msg("Error creating like")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error liking thread",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Thread liked",
		"error":   nil,
	})
}

// Handle unlike thread request
func UnlikeThreadHandler(c *gin.Context) {
	username := services.JwtHandler.GetUsernameFromJwt(c)
	if username == "" {
		return
	}

	threadID := c.Param("thread_id")
	utils.Logger.Trace().Str("thread_id", threadID).Msg("")

	if !services.Threads.ThreadExists(threadID) {
		utils.Logger.Error().Msg("Thread does not exist")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Thread does not exist",
			"error":   nil,
		})
	}

	hasLiked := services.Likes.HasLiked(username, threadID)
	if !hasLiked {
		utils.Logger.Trace().Msg("User has not liked thread")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User has not liked thread",
			"error":   nil,
		})
		return
	}

	err := services.Likes.RemoveLike(username, threadID)
	if err != nil {
		utils.Logger.Err(err).Msg("Error removing like")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error unliking thread",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Thread unliked",
		"error":   nil,
	})
}
