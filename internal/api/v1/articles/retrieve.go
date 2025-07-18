package articles

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve all articles
func GetAllArticlesHandler(c *gin.Context) {
	// Get user uid from jwt
	uid := services.JwtHandler.GetUidFromJwt(c)

	// Retrieve query params
	size, err := strconv.Atoi(c.DefaultQuery("size", constants.DEFAULT_PAGE_SIZE))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	desc := c.DefaultQuery("desc", constants.DEFAULT_SORT_DESCENDING) == "true"
	sort := c.DefaultQuery("sort", constants.DEFAULT_SORT_COLUMN)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	utils.Logger.Debug().Int("size", size).Bool("desc", desc).Str("sort", sort).Int("page", page).Msg("Get all articles request received.")

	articles, err := services.Articles.GetArticles(sort, size, desc, page, uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting articles")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error retrieving articles" + err.Error(),
		})
		return
	}

	metadata, err := services.Articles.GetMetadata()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered retrieving metadata")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "error encountered retrieving metadata" + err.Error(),
		})
		return
	}

	// Set number of pages
	metadata.NumPages = (metadata.Total + size - 1) / size

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"articles": articles,
		"metadata": metadata,
	})
}

// Retrieve single article based on article_id parameter
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
