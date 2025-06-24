package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateArticleHandler(c *gin.Context) {
	var createArticleRequest CreateArticleRequest

	// Bind with form
	if err := c.ShouldBind(&createArticleRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new article request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Malformed request",
		})
		return
	}

	// Get author uid from JWT token
	author := services.JwtHandler.GetUidFromJwt(c)
	utils.Logger.Info().Msg("New article request received from " + author)

	id, err := services.Articles.CreateNewArticle(author, createArticleRequest.Title, createArticleRequest.Content)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating new article " + err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"message":   "Article created successfully",
		"articleId": id,
	})
}
