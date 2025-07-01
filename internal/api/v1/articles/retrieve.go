package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func GetOneArticleHandler(c *gin.Context) {
	articleID := c.Param("article_id")

	// Get uid from jwt
	uid := services.JwtHandler.GetUidFromJwt(c)

	article, comments, err := services.Articles.GetArticle(articleID, uid)

	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting article")

		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "article not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"article":  article,
		"comments": comments,
	})
}
