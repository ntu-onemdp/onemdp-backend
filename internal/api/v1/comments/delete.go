package comments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("comment_id")

	utils.Logger.Trace().Str("commentID", commentID).Msgf("Delete comment request received for comment %s", commentID)

	// Get author from JWT token
	author := services.JwtHandler.GetUidFromJwt(c)
	utils.Logger.Info().Msgf("Delete commnent request received from %s for comment %s", author, commentID)

	// Delete comment
	if err := services.Comments.Delete(commentID, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting comment")

		if (err == utils.ErrUnauthorized{}) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized to delete comment. You need to be a staff/admin or the original author to delete the comment",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error deleting comment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment successfully deleted",
	})
}
