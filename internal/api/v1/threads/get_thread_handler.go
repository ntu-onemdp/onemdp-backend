package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type GetThreadHandler struct {
	ThreadService *services.ThreadService
}

func (h *GetThreadHandler) HandleGetThread(c *gin.Context) {
	threadId := c.Param("thread_id")

	thread, posts, err := h.ThreadService.GetThread(threadId)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting thread")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"thread":  thread,
		"posts":   posts,
	})
}
