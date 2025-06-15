package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// For now this handler will only update one user role per request
type UpdateUserRoleRequest struct {
	Uid  string `json:"uid" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func UpdateRoleHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Update user role request received")

	// Parse request
	var updateUserRoleRequest UpdateUserRoleRequest
	if err := c.BindJSON(&updateUserRoleRequest); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding request to UpdateUsersRoleRequest")
		c.JSON(400, nil)
		return
	}

	utils.Logger.Trace().Interface("updateUserRoleRequest", updateUserRoleRequest).Msg("Parsed request")

	err := services.Users.UpdateRole(updateUserRoleRequest.Uid, updateUserRoleRequest.Role)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when updating user role")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	// Success response: return 200
	c.JSON(http.StatusOK, nil)
}
