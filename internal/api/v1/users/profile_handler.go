package users

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ProfileHandler struct {
	UserService *services.UserService
}

type HasPasswordChangedResponse struct {
	Username         string `json:"username"`
	Password_changed bool   `json:"password_changed"`
}

func (h *ProfileHandler) HandleHasPasswordChanged(c *gin.Context) {
	utils.Logger.Trace().Msg("Password changed query received")

	// Retrieve username and jwt
	username := c.Param("username")
	tokenString := c.Request.Header.Get("Authorization")

	// Remove "Bearer " prefix if included
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// JWT does not match username
	if !utils.ValidateUsername(username, tokenString) {
		c.JSON(401, nil)
		return
	}

	password_changed, err := h.UserService.HasPasswordChanged(username)
	if err != nil {
		utils.Logger.Error().Err(err)
		c.JSON(500, nil)
	}

	c.JSON(200, HasPasswordChangedResponse{
		Username:         username,
		Password_changed: password_changed,
	})
}
