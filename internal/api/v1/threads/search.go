package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchThreadsHandler(c *gin.Context) {
	params := c.Query("params")

	c.JSON(http.StatusOK, gin.H{
		"message": "This feature has not been implemented yet",
		"search":  params,
	})
}
