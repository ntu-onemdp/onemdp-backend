package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func UpdatePostHandler(c *gin.Context) {
	// Bind with post object
	var updatedPost models.Post

	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Malformed request",
		})
		return
	}

	utils.Logger.Debug().Interface("updatedPost", updatedPost).Msg("Update post request")

	// Safeguard against post id manipulation
	if updatedPost.PostID != c.Param("post_id") {
		utils.Logger.Error().Msg("Post ID mismatch")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Post ID mismatch",
		})
		return
	}

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := services.JwtHandler.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	author := claim.Uid
	utils.Logger.Info().Msg("Update post request received from " + author)

	// Update post
	err = services.Posts.UpdatePost(updatedPost, claim)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating post")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error updating post",
		})
		return
	}

	// Update thread last activity
	if updatedPost.IsHeader {
		// Header post: update title, preview, and last activity
		err = services.Threads.UpdateThread(updatedPost.ThreadId, updatedPost.Title, updatedPost.PostContent, claim)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error updating thread")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error updating thread",
			})
			return
		}
	} else {
		// All other posts: update only last activity
		err = services.Threads.UpdateThreadLastActivity(updatedPost.ThreadId)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error updating thread last activity")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error updating thread last activity",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post updated successfully",
	})
}
