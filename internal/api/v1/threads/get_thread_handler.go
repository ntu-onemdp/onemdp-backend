package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type GetThreadHandler struct {
	threadService *services.ThreadService
	likeService   *services.LikeService
}

func NewGetThreadHandler(threadService *services.ThreadService, likeService *services.LikeService) *GetThreadHandler {
	return &GetThreadHandler{
		threadService: threadService,
		likeService:   likeService,
	}
}

func (h *GetThreadHandler) HandleGetThread(c *gin.Context) {
	threadId := c.Param("thread_id")

	thread, posts, err := h.threadService.GetThread(threadId)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting thread")

		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "thread not found",
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
		"success":     true,
		"thread":      thread,
		"posts":       posts,
		"num_replies": len(posts) - 1,
	})
}
