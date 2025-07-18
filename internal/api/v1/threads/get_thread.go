package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Retrieve all threads in page
func GetAllThreadsHandler(c *gin.Context) {
	// Get uid from jwt
	uid := services.JwtHandler.GetUidFromJwt(c)

	size := c.GetInt("size")
	if size == 0 {
		size = constants.DEFAULT_PAGE_SIZE
	}
	desc := c.DefaultQuery("desc", constants.DEFAULT_SORT_DESCENDING) == "true"
	sort := c.DefaultQuery("sort", constants.DEFAULT_SORT_COLUMN)
	page := c.GetInt("page")

	// Page not provided: set to first page
	if page == 0 {
		page = 1 // First page
	}

	utils.Logger.Debug().Int("size", size).Bool("desc", desc).Str("sort", sort).Int("page", page).Msg("Get all threads request received.")

	threads, err := services.Threads.GetThreads(sort, size, desc, page, uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting threads")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error getting threads",
			"error":   err.Error(),
		})
		return
	}

	metadata, err := services.Threads.GetThreadsMetadata()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting threads metadata")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error getting threads metadata",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"threads":     threads,
		"num_threads": metadata.NumThreads,
		"num_pages":   (metadata.NumThreads + size - 1) / size,
	})
}

// Retrieve individual thread
func GetOneThreadHandler(c *gin.Context) {
	threadId := c.Param("thread_id")

	// Get uid from jwt
	uid := services.JwtHandler.GetUidFromJwt(c)

	thread, posts, err := services.Threads.GetThread(threadId, uid)
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
