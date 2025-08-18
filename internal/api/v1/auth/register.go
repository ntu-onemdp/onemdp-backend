package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/semester"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Register user request from frontend
type req struct {
	user
	Code string `json:"code" binding:"required"` // Enrolment code
}

func RegisterUserHandler(c *gin.Context) {
	utils.Logger.Info().Msg("Register user request received.")
	var req req

	// Bind with form
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding register user request with struct.")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error registering user. Check that body is correct",
			"error":   err,
		})
		return
	}

	// Check if code matches
	if req.Code != semester.Service.GetCode() {
		utils.Logger.Warn().Str("Provided code", req.Code).Str("Registration code", semester.Service.GetCode()).Msg("Registration code does not match.")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Enrolment code provided is incorrect",
		})
		return
	}

	utils.Logger.Trace().Msg("Code provided is correct. Registering new user.")

	// Register user
	if err := services.Users.RegisterUser(req.Uid, req.UserMetadata.Email, req.UserMetadata.Name); err != nil {
		utils.Logger.Error().Err(err).Msg("Error encountered when registering user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error encountered when registering user",
			"error":   err,
		})
		return
	}

	utils.Logger.Info().Str("uid", req.Uid).Msgf("User %s successfully registered via enrolment code", req.UserMetadata.Name)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User successfully registered into the system.",
	})
}
