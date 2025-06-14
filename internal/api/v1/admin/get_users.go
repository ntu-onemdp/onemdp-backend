package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Get all users
func GetAllUsersHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Get users request received")

	users, err := services.Users.GetAllUsersAdmin()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when getting all users")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

// Get individual user
func GetOneUserHandler(c *gin.Context) {
	uid := c.Param("uid")
	utils.Logger.Info().Msg("Get user request received for " + uid)

	user, err := services.Users.GetUserAdmin(uid)
	if user == nil {
		utils.Logger.Error().Msg("User not found")
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when getting user")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
