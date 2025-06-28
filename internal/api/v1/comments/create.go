package comments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// New comment request from the frontend
type CreateCommentRequest struct {
	ArticleID string `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

func CreateCommentHandler(c *gin.Context) {
	var request CreateCommentRequest

	// Bind with form
	if err := c.ShouldBind(&request); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new comment request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Malformed request. " + err.Error(),
		})
		return
	}

	// Get author uid from jwt token
	author := services.JwtHandler.GetUidFromJwt(c)
	utils.Logger.Info().Msg("New comment request received from " + author)

	id, err := services.Comments.Create(author, request.ArticleID, request.Content)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating new comment " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "Comment created successfully",
		"comment_id": id,
	})
}
