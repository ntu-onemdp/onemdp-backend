package karma

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/karma"
)

func UpdateKarmaHandler(c *gin.Context) {
	var req karma.Karma
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Error binding karma request",
		})
		return
	}

	if err := karma.Service.UpdateSettings(karma.Karma{
		CreateThreadPts:  req.CreateThreadPts,
		CreateArticlePts: req.CreateArticlePts,
		CreateCommentPts: req.CreateCommentPts,
		CreatePostPts:    req.CreatePostPts,
		LikePts:          req.LikePts,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Error updating karma settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Karma settings updated successfully",
	})
}
