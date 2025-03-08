package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type NewPostHandler struct {
	PostService *services.PostService
}

// Frontend request to create a new post. Get author of post from JWT token.
// This request uses JSON instead of form data.
type NewPostRequest struct {
	Title    string    `json:"title" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	ReplyTo  string    `json:"reply_to" binding:"required"` // NA if not a reply
	ThreadId uuid.UUID `json:"thread_id" binding:"required"`
}

func (h *NewPostHandler) HandleNewPost(c *gin.Context) {
	var newPostRequest NewPostRequest

	// Bind with JSON
	if err := c.ShouldBindJSON(&newPostRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing new post request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success":  false,
			"errorMsg": "Malformed request",
		})
		return
	}

	// For debugging purposes
	utils.Logger.Debug().Interface("newPostRequest", newPostRequest).Msg("New post request")

	// Get author from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := utils.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	author := claim.Username
	utils.Logger.Info().Msg("New post request received from " + author)

	// Check if reply to is NA
	var replyTo *string
	if newPostRequest.ReplyTo == "NA" {
		replyTo = nil
	} else {
		replyTo = &newPostRequest.ReplyTo
	}

	// Create new post
	err = h.PostService.CreateNewPost(author, replyTo, newPostRequest.ThreadId, newPostRequest.Title, newPostRequest.Content)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating new post")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Post created successfully",
	})
}
