package threads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type NewThreadHandler struct {
	ThreadService *services.ThreadService
}

// Frontend request to create a new thread. Get author from JWT token
type NewThreadRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}

func (h *NewThreadHandler) HandleNewThread(c *gin.Context) {
	var newThreadRequest NewThreadRequest

	// Bind with form
	if err := c.ShouldBind(&newThreadRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new thread request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"errorMsg": "Malformed request",
		})
		return
	}

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := utils.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	author := claim.Username
	utils.Logger.Info().Msg("New thread request received from " + author)

	// Create new thread
	err = h.ThreadService.CreateNewThread(author, newThreadRequest.Title, newThreadRequest.Content)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating new thread")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Thread created successfully",
	})
}
