package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Update validation status of a post
func ValidatePostHandler(c *gin.Context) {
	// Get uid from JWT token
	uid := services.JwtHandler.GetUidFromJwt(c)

	postID := c.Param("post_id")

	// Retrieve new validation status from request body
	var reqBody struct {
		ValidationStatus string `json:"validation_status" binding:"required,oneof=unverified validated refuted"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	utils.Logger.Info().Str("uid", uid).Str("post_id", postID).Str("validation_status", reqBody.ValidationStatus).Msg("Validate post request received")

	err := services.Posts.UpdateValidationStatus(postID, reqBody.ValidationStatus, uid)
	if err != nil {
		utils.Logger.Error().Err(err).Str("uid", uid).Str("post_id", postID).Str("validation_status", reqBody.ValidationStatus).Msg("Error updating validation status of post")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error updating validation status of post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post validation status updated successfully",
	})
}
