package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Get all users
func GetAllUsersHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Get users request received")

	users, err := services.Users.GetAllUsersInformation()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when getting all users")
		c.JSON(500, nil)
		return
	}

	// Success response: return 200
	c.JSON(200, users)
}

// Get individual user
func GetOneUserHandler(c *gin.Context) {
	username := c.Param("username")
	utils.Logger.Info().Msg("Get user request received for " + username)

	user, err := services.Users.GetUserInformation(username)
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
