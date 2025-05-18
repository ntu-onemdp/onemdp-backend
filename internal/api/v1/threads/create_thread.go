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
}

func CreateThreadHandler(c *gin.Context) {
	var createThreadRequest CreateThreadRequest

	// Bind with form
	if err := c.ShouldBind(&createThreadRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new thread request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Malformed request",
		})
		return
	}

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := services.JwtHandler.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	author := claim.Username
	utils.Logger.Info().Msg("New thread request received from " + author)

	id, err := services.Threads.CreateNewThread(author, createThreadRequest.Title, createThreadRequest.Content)
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
