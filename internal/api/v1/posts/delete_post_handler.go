package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type DeletePostHandler struct {
	PostService *services.PostService
}

func (h *DeletePostHandler) HandleDeletePost(c *gin.Context) {
	postID := c.Param("post_id")

	// For debugging purposes
	utils.Logger.Info().Str("postID", postID).Msg("Delete post request received")

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := utils.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	utils.Logger.Info().Msg("Delete post request received from " + claim.Username)

	// Delete post
	err = h.PostService.DeletePost(postID, claim)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting post")

		if (err == utils.ErrUnauthorized{}) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":  false,
				"errorMsg": "Unauthorized to delete post. You need to be a staff/admin or the original author to delete the post",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":  false,
			"errorMsg": "Error deleting post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post successfully deleted",
	})
}
