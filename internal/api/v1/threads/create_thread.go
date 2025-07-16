package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Frontend request to create a new thread. Get author from JWT token
type CreateThreadRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	IsAnon  bool   `json:"is_anon" binding:"required"`
}

func CreateThreadHandler(c *gin.Context) {
	var createThreadRequest CreateThreadRequest

	utils.Logger.Debug().Interface("body", c.Request.Body).Msg("request body")

	// Bind with form
	if err := c.ShouldBindJSON(&createThreadRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new thread request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Malformed request",
		})
		return
	}

	// Get author uid from JWT token
	author := services.JwtHandler.GetUidFromJwt(c)
	utils.Logger.Info().Msg("New thread request received from " + author)

	id, err := services.Threads.CreateNewThread(author, createThreadRequest.Title, createThreadRequest.Content, createThreadRequest.IsAnon)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating new thread " + err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  "Thread created successfully",
		"threadId": id,
	})
}
