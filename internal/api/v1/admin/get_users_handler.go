package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type GetUsersHandler struct {
	UserService *services.UserService
}

// Get all users
func (h *GetUsersHandler) HandleGetUsers(c *gin.Context) {
	utils.Logger.Info().Msg("Get users request received")

	users, err := h.UserService.GetAllUsersInformation()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when getting all users")
		c.JSON(500, nil)
		return
	}

	// Success response: return 200
	c.JSON(200, users)
}

// Get individual user
func (h *GetUsersHandler) HandleGetUser(c *gin.Context) {
	username := c.Param("username")
	utils.Logger.Info().Msg("Get user request received for " + username)

	user, err := h.UserService.GetUserInformation(username)
	if user == nil {
		utils.Logger.Error().Msg("User not found")
		c.JSON(404, nil)
		return
	}
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when getting user")
		c.JSON(500, nil)
		return
	}

	// Success response: return 200
	c.JSON(200, user)
}
