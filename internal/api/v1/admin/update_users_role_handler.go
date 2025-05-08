package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type UpdateUsersRoleHandler struct {
	AuthService *services.AuthService
}

// For now this handler will only update one user role per request
type UpdateUserRoleRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (h *UpdateUsersRoleHandler) HandleUpdateUsersRole(c *gin.Context) {
	utils.Logger.Info().Msg("Update user role request received")

	// Parse request
	var updateUserRoleRequest UpdateUserRoleRequest
	if err := c.BindJSON(&updateUserRoleRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding request to UpdateUsersRoleRequest")
		c.JSON(400, nil)
		return
	}

	err := h.AuthService.UpdateRole(updateUserRoleRequest.Username, updateUserRoleRequest.Role)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when updating user role")
		c.JSON(500, nil)
		return
	}

	// Success response: return 200
	c.JSON(200, nil)
}
