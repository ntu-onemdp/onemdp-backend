package threads

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
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

// Retrieve all threads in page
func (h *GetThreadHandler) HandleGetThreads(c *gin.Context) {
	const DEFAULT_PAGE_SIZE = 25
	const DEFAULT_SORT_COLUMN = models.TIME_CREATED_COL
	const DEFAULT_SORT_DESCENDING = true

	size := c.GetInt("size")
	if size == 0 {
		size = DEFAULT_PAGE_SIZE
	}
	desc := c.DefaultQuery("desc", strconv.FormatBool(DEFAULT_SORT_DESCENDING)) == "true"
	sort := c.DefaultQuery("sort", string(DEFAULT_SORT_COLUMN)) // Defaults to time_created if column name is invalid
	timestamp := c.GetTime("timestamp")
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	utils.Logger.Debug().Int("size", size).Bool("desc", desc).Str("sort", sort).Time("timestamp", timestamp).Msg("")

	threads, err := h.threadService.GetThreads(sort, size, desc, timestamp)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting threads")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error getting threads",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"threads": threads,
	})
}

// Retrieve individual thread
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
