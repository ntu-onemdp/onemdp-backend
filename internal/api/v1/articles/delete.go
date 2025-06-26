package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func DeleteArticleHandler(c *gin.Context) {
	articleId := c.Param("article_id")

	// Get author from JWT token
	author := services.JwtHandler.GetUidFromJwt(c)

	utils.Logger.Info().Msg("Delete article request received from " + author)

	err := services.Articles.DeleteArticle(articleId, author)
	if err == utils.NewErrUnauthorized() {
		utils.Logger.Error().Err(err).Msg("User is student and not author. Unauthorized to delete article")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized to delete article. You need to be a staff/admin or the original author to delete the article",
		})
		return
	} else if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting article")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Error deleting article: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Article deleted successfully",
	})
}
