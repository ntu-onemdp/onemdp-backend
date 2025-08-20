package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// uid + user object retrieved from supabase
type user struct {
	Uid          string       `json:"uid" binding:"required"`
	UserMetadata userMetadata `json:"user_metadata" binding:"required"`
}

type userMetadata struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"full_name" binding:"required"`
}

type LoginResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Jwt     *string             `json:"jwt"`
	User    *models.UserProfile `json:"user"`
}

// Implemented for SSO. After SSO login, the handler will return the JWT and user profile
func LoginHandler(c *gin.Context) {
	utils.Logger.Trace().Msg("Login request received")
	var user user

	// Bind with form
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing login request")
		response := LoginResponse{
			Success: false,
			Error:   "Malformed request",
		}

		c.JSON(http.StatusBadRequest, &response)
		return
	}

	utils.Logger.Debug().Str("uid", user.Uid).Interface("user metadata", user.UserMetadata).Msg("Login request binded")

	// Check if user is pending registration
	isPending, err := services.Users.IsUserPending(user.UserMetadata.Email)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking if user is pending")
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	if isPending {
		utils.Logger.Debug().Msg("User is pending registration")

		if err := services.Users.RegisterUserFromPending(user.Uid, user.UserMetadata.Email, user.UserMetadata.Name); err != nil {
			utils.Logger.Error().Err(err).Msg("Error registering user")
			response := LoginResponse{
				Success: false,
				Error:   "Internal server error",
			}
			c.JSON(http.StatusInternalServerError, &response)
			return
		}
	}

	// Return user profile
	profile, err := services.Users.GetProfile(user.Uid)
	if err != nil {
		utils.Logger.Debug().Msg("User profile not found")
		response := LoginResponse{
			Success: false,
			Error:   "User not registered",
		}
		c.JSON(http.StatusNotFound, &response) // 16/08/2025: Changed from StatusUnauthorized to StatusNotFound.
		return
	}

	// Generate jwt
	claim := models.NewClaim(user.Uid)
	jwt, err := services.JwtHandler.GenerateJwt(claim)
	if err != nil {
		utils.Logger.Error().Err(err)
		c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	response := LoginResponse{
		Success: true,
		Jwt:     &jwt,
		User:    profile,
	}
	c.JSON(http.StatusOK, &response)
}
