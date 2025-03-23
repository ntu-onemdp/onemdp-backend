package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ChangePasswordHandler struct {
	AuthService *services.AuthService
}

// Password change form sent from frontend.
type ChangePasswordForm struct {
	OldPassword string `form:"old_password" binding:"required"` // Plaintext password
	NewPassword string `form:"new_password" binding:"required"` // Plaintext password
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *ChangePasswordHandler) HandleChangeUserPassword(c *gin.Context) {
	var form ChangePasswordForm
	username := c.Param("username")
	tokenString := c.Request.Header.Get("Authorization")

	utils.Logger.Info().Msg("Change password request received for " + username)

	// Bind with form
	if err := c.ShouldBind(&form); err != nil {
		response := ChangePasswordResponse{
			Success: false,
			Message: "Malformed request",
		}

		c.JSON(400, &response)
		return
	}

	// Validate JWT
	if !utils.JwtHandler.ValidateUsername(username, tokenString) {
		response := ChangePasswordResponse{
			Success: false,
			Message: "Error: Invalid JWT",
		}

		c.JSON(401, &response)
		return
	}

	// Check if old password matches old password in database
	isAuthenticated, _, _ := h.AuthService.AuthenticateUser(username, form.OldPassword)
	if !isAuthenticated {
		response := ChangePasswordResponse{
			Success: false,
			Message: "Error: Incorrect old password",
		}

		c.JSON(200, &response)
		return
	}

	// Change password
	err := h.AuthService.UpdateUserPassword(username, form.NewPassword)
	if err != nil {
		response := ChangePasswordResponse{
			Success: false,
			Message: err.Error(),
		}

		c.JSON(200, &response)
		return
	}

	response := ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}

	c.JSON(200, &response)
}
