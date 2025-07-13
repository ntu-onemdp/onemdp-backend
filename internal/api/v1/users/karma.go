package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

func GetRankingsHandler(c *gin.Context) {
	semester := c.Query("semester")

	users, err := services.Users.GetTopKarma(semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
