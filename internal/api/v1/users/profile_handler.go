package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ProfileHandler struct {
	UserService *services.UserService
}

type HasPasswordChangedResponse struct {
	Username         string `json:"username"`
	Password_changed bool   `json:"password_changed"` // False if user has not changed his password yet
}

func (h *ProfileHandler) HandleHasPasswordChanged(c *gin.Context) {
	// Retrieve username and jwt
	username := c.Param("username")
	tokenString := c.Request.Header.Get("Authorization")

	utils.Logger.Info().Msg(fmt.Sprintf("Password changed query received for %s", username))

	// JWT does not match username
	if !utils.ValidateUsername(username, tokenString) {
		c.JSON(401, nil)
		return
	}

	password_changed, err := h.UserService.HasPasswordChanged(username)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(500, nil)
	}

	c.JSON(200, HasPasswordChangedResponse{
		Username:         username,
		Password_changed: password_changed,
	})
}

// Retrieve public user profile information as defined in models.UserProfile
func (h *ProfileHandler) HandleGetUserProfile(c *gin.Context) {
	username := c.Param("username")

	utils.Logger.Info().Str("username", username).Msg(fmt.Sprintf("Get user profile request received for %s", username))

	profile, err := h.UserService.GetUserProfile(username)
	if err != nil {
		utils.Logger.Debug().Msg("profile may be nil. returning 404 here")
		utils.Logger.Error().Err(err).Msg("")
		c.JSON(404, nil)
		return
	}

	c.JSON(200, profile)
}
